package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	_dictionaryGrpc "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc"
	pb "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc/gen"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lestrrat-go/jwx/jwk"
	_dictionaryHttpDelivery "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/http"
	_dictionaryRepo "github.com/vovancho/lingua-cat-go/dictionary/dictionary/repository/postgres"
	_dictionaryUcase "github.com/vovancho/lingua-cat-go/dictionary/dictionary/usecase"

	_ "github.com/lib/pq"
)

var cachedJWK jwk.Set

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file, using default config")
	}

	var err error
	cachedJWK, err = loadJWKFromPEM("./internal/misc/public.pem")
	if err != nil {
		slog.Error("Не удалось загрузить JWK: %v", err)
	} else {
		slog.Info("JWK загружен")
	}
}

func main() {
	dbDsn := os.Getenv("DB_DSN")
	dbConn := initDbConn(dbDsn)

	defer func() {
		err := dbConn.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	validate, trans, err := initValidator()
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	dictionaryRepo := _dictionaryRepo.NewPostgresDictionaryRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	du := _dictionaryUcase.NewDictionaryUseCase(dictionaryRepo, validate, timeoutContext)
	_dictionaryHttpDelivery.NewDictionaryHandler(router, validate, du)

	server := http.Server{
		Addr:    ":80",
		Handler: response.ErrorMiddleware(AuthMiddleware(router), trans),
	}

	// gRPC-сервер
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor))
	dictionaryHandler := _dictionaryGrpc.NewDictionaryHandler(validate, du)
	pb.RegisterDictionaryServiceServer(grpcServer, dictionaryHandler)
	// Запуск gRPC-сервера в отдельной горутине
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			slog.Error("Failed to listen for gRPC", "error", err)
			panic(err)
		}
		slog.Info("gRPC server is listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("Failed to serve gRPC", "error", err)
			panic(err)
		}
	}()

	fmt.Println("Server is listening on port 80")
	server.ListenAndServe()
}

func initDbConn(dsn string) *sqlx.DB {
	dbConn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic("не удалось сконфигурировать подключение к БД")
	}

	if err = dbConn.Ping(); err != nil {
		panic("не удалось подключиться к серверу БД")
	}

	return dbConn
}

func initValidator() (*validator.Validate, ut.Translator, error) {
	validate := validator.New()
	uni := ut.New(ru.New(), ru.New())
	trans, _ := uni.GetTranslator("ru")

	if err := rutranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	if err := domain.RegisterAll(validate, trans); err != nil {
		return nil, nil, fmt.Errorf("failed to register domain validations: %w", err)
	}

	return validate, trans, nil
}

func loadJWKFromPEM(pemFile string) (jwk.Set, error) {
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA, got %T", cert.PublicKey)
	}

	key := jwk.NewRSAPublicKey()
	if err := key.FromRaw(rsaPubKey); err != nil {
		return nil, err
	}

	set := jwk.NewSet()
	set.Add(key)
	return set, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
			return
		}

		key, ok := cachedJWK.Get(0)
		var rsaPubKey rsa.PublicKey
		key.Raw(&rsaPubKey)

		token, err := jwt.Parse(
			[]byte(tokenStr),
			jwt.WithVerify(jwa.RS256, rsaPubKey),
			jwt.WithValidate(true),
		)
		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		sub, ok := token.Get("sub")
		if !ok {
			http.Error(w, "Отсутствует sub в токене", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Отсутствуют метаданные")
	}

	authHeaders, ok := md["authorization"]
	if !ok || len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Отсутствует заголовок Authorization")
	}

	tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if tokenStr == authHeaders[0] {
		return nil, status.Error(codes.Unauthenticated, "Неверный формат заголовка Authorization")
	}

	key, ok := cachedJWK.Get(0)
	var rsaPubKey rsa.PublicKey
	key.Raw(&rsaPubKey)

	token, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithVerify(jwa.RS256, rsaPubKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Неверный токен")
	}

	sub, ok := token.Get("sub")
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Отсутствует sub в токене")
	}

	ctx = context.WithValue(ctx, "userID", sub)
	return handler(ctx, req)
}
