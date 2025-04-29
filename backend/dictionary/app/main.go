package main

import (
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	_dictionaryGrpc "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc"
	pb "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc/gen"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_dictionaryHttpDelivery "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/http"
	_dictionaryRepo "github.com/vovancho/lingua-cat-go/dictionary/dictionary/repository/postgres"
	_dictionaryUcase "github.com/vovancho/lingua-cat-go/dictionary/dictionary/usecase"

	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file, using default config")
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
		Handler: response.ErrorMiddleware(router, trans),
	}

	// gRPC-сервер
	grpcServer := grpc.NewServer()
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
