//go:build wireinject
// +build wireinject

package wire

import (
	"net/http"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	_internalHttp "github.com/vovancho/lingua-cat-go/analytics/analytics/delivery/http"
	"github.com/vovancho/lingua-cat-go/analytics/analytics/repository/clickhouse"
	"github.com/vovancho/lingua-cat-go/analytics/analytics/usecase"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
	"github.com/vovancho/lingua-cat-go/analytics/internal/config"
	"github.com/vovancho/lingua-cat-go/analytics/internal/db"
	"github.com/vovancho/lingua-cat-go/analytics/internal/response"
	"github.com/vovancho/lingua-cat-go/analytics/internal/translator"
	_internalValidator "github.com/vovancho/lingua-cat-go/analytics/internal/validator"
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
		getClickHouseDB,

		// Переводчик
		translator.NewTranslator,

		// Валидатор
		_internalValidator.NewValidator,

		// Аутентификация
		ProvidePublicKeyPath,
		auth.NewAuthService,

		//// Репозиторий
		clickhouse.NewClickhouseExerciseCompleteRepository,

		//// Use case
		ProvideUseCaseTimeout,
		usecase.NewExerciseCompleteUseCase,

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
func getClickHouseDB(db *sqlx.DB) db.DB {
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
	exerciseCompleteUcase domain.ExerciseCompleteUseCase,
) *http.Server {
	router := http.NewServeMux()
	_internalHttp.NewExerciseCompleteHandler(router, validate, authService, exerciseCompleteUcase)
	return &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: response.ErrorMiddleware(authService.AuthMiddleware(router), trans),
	}
}

//// ProvideKafkaExerciseCompletedTopic возвращает имя топика о выполненных упражнениях
//func ProvideKafkaExerciseCompletedTopic(cfg *config.Config) usecase.ExerciseCompletedTopic {
//	return usecase.ExerciseCompletedTopic(cfg.KafkaExerciseCompletedTopic)
//}
