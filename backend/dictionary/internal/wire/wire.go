//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_internalGrpc "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc"
	pb "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc/gen"
	_internalHttp "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/http"
	"github.com/vovancho/lingua-cat-go/dictionary/dictionary/repository/postgres"
	"github.com/vovancho/lingua-cat-go/dictionary/dictionary/usecase"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/auth"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/config"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/db"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/translator"
	_internalValidator "github.com/vovancho/lingua-cat-go/dictionary/internal/validator"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

// App представляет приложение с конфигурацией и серверами
type App struct {
	Config     *config.Config
	HTTPServer *http.Server
	GRPCServer *grpc.Server
	DB         *sqlx.DB
}

// NewApp создаёт новый экземпляр App
func NewApp(cfg *config.Config, httpServer *http.Server, grpcServer *grpc.Server, db *sqlx.DB) *App {
	return &App{
		Config:     cfg,
		HTTPServer: httpServer,
		GRPCServer: grpcServer,
		DB:         db,
	}
}

func InitializeApp() (*App, error) {
	wire.Build(
		// Конфигурация
		config.Load,

		// База данных
		ProvideDSN,
		db.NewDB,
		getPostgresDB,

		// Переводчик
		translator.NewTranslator,

		// Валидатор
		_internalValidator.NewValidator,

		// Аутентификация
		ProvidePublicKeyPath,
		auth.NewAuthService,

		// Репозиторий
		postgres.NewPostgresDictionaryRepository,

		// Use case
		ProvideUseCaseTimeout,
		usecase.NewDictionaryUseCase,

		// HTTP Delivery
		newHTTPServer,

		// gRPC Delivery
		newGRPCServer,

		// App
		NewApp,
	)
	return &App{}, nil
}

func ProvideDSN(cfg *config.Config) db.DSN {
	return db.DSN(cfg.DBDSN)
}

func ProvidePublicKeyPath(cfg *config.Config) auth.PublicKeyPath {
	return auth.PublicKeyPath(cfg.AuthPublicKeyPath)
}

// getPostgresDB возвращает *sqlx.DB как db.DB
func getPostgresDB(db *sqlx.DB) db.DB {
	return db
}

// getUseCaseTimeout возвращает таймаут для use case из конфигурации
func ProvideUseCaseTimeout(cfg *config.Config) usecase.Timeout {
	return usecase.Timeout(time.Duration(cfg.Timeout) * time.Second)
}

// newHTTPServer создаёт новый HTTP-сервер
func newHTTPServer(
	cfg *config.Config,
	validate *validator.Validate,
	trans ut.Translator,
	authService *auth.AuthService,
	dictionaryUcase domain.DictionaryUseCase,
) *http.Server {
	router := http.NewServeMux()
	_internalHttp.NewDictionaryHandler(router, validate, dictionaryUcase)

	// Register gRPC-Gateway handlers
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterDictionaryServiceHandlerFromEndpoint(context.Background(), gwmux, cfg.GRPCPort, opts); err != nil {
		panic(err)
	}

	mainMux := http.NewServeMux()
	mainMux.Handle("/grpc-gw-swagger.json", http.FileServer(http.Dir("doc")))
	mainMux.Handle("/grpc-gateway/", gwmux)

	mainMux.Handle("/swagger.json", http.FileServer(http.Dir("doc")))
	mainMux.Handle("/", response.ErrorMiddleware(authService.AuthMiddleware(router), trans))

	return &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mainMux,
	}
}

// newGRPCServer создаёт новый gRPC-сервер
func newGRPCServer(
	validate *validator.Validate,
	authService *auth.AuthService,
	dictionaryUcase domain.DictionaryUseCase,
) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authService.AuthInterceptor))
	pb.RegisterDictionaryServiceServer(grpcServer, _internalGrpc.NewDictionaryHandler(validate, dictionaryUcase))
	return grpcServer
}
