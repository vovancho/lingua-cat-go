//go:build wireinject
// +build wireinject

package wire

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	_internalHttp "github.com/vovancho/lingua-cat-go/exercise/exercise/delivery/http"
	_internalGrpc "github.com/vovancho/lingua-cat-go/exercise/exercise/repository/grpc"
	"github.com/vovancho/lingua-cat-go/exercise/exercise/repository/postgres"
	"github.com/vovancho/lingua-cat-go/exercise/exercise/usecase"
	"github.com/vovancho/lingua-cat-go/exercise/internal/auth"
	"github.com/vovancho/lingua-cat-go/exercise/internal/config"
	"github.com/vovancho/lingua-cat-go/exercise/internal/db"
	"github.com/vovancho/lingua-cat-go/exercise/internal/response"
	"github.com/vovancho/lingua-cat-go/exercise/internal/translator"
	_internalValidator "github.com/vovancho/lingua-cat-go/exercise/internal/validator"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

// App представляет приложение с конфигурацией и серверами
type App struct {
	Config     *config.Config
	HTTPServer *http.Server
	DB         *sqlx.DB
}

// NewApp создаёт новый экземпляр App
func NewApp(cfg *config.Config, httpServer *http.Server, db *sqlx.DB) *App {
	return &App{
		Config:     cfg,
		HTTPServer: httpServer,
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

		// GRPC соединение
		ProvideGRPCConn,

		// Репозиторий
		postgres.NewPostgresExerciseRepository,
		postgres.NewPostgresTaskRepository,
		_internalGrpc.NewGrpcDictionaryRepository,

		// Use case
		ProvideUseCaseTimeout,
		usecase.NewExerciseUseCase,
		usecase.NewTaskUseCase,
		usecase.NewDictionaryUseCase,

		// HTTP Delivery
		newHTTPServer,

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

func ProvideGRPCConn(cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cfg.DictionaryGRPCAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// newHTTPServer создаёт новый HTTP-сервер
func newHTTPServer(
	cfg *config.Config,
	validate *validator.Validate,
	trans ut.Translator,
	authService *auth.AuthService,
	exerciseUcase domain.ExerciseUseCase,
	taskUcase domain.TaskUseCase,
) *http.Server {
	router := http.NewServeMux()
	_internalHttp.NewExerciseHandler(router, validate, authService, exerciseUcase)
	_internalHttp.NewTaskHandler(router, validate, authService, taskUcase, exerciseUcase)
	return &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: response.ErrorMiddleware(authService.AuthMiddleware(router), trans),
	}
}
