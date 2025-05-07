//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	_internalHttp "github.com/vovancho/lingua-cat-go/analytics/analytics/delivery/http"
	_internalKafka "github.com/vovancho/lingua-cat-go/analytics/analytics/delivery/kafka"
	"github.com/vovancho/lingua-cat-go/analytics/analytics/repository/clickhouse"
	"github.com/vovancho/lingua-cat-go/analytics/analytics/usecase"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
	"github.com/vovancho/lingua-cat-go/analytics/internal/config"
	"github.com/vovancho/lingua-cat-go/analytics/internal/db"
	"github.com/vovancho/lingua-cat-go/analytics/internal/response"
	"github.com/vovancho/lingua-cat-go/analytics/internal/translator"
	_internalValidator "github.com/vovancho/lingua-cat-go/analytics/internal/validator"
	"net/http"
	"time"
)

// App представляет приложение с конфигурацией и серверами
type App struct {
	Config           *config.Config
	HTTPServer       *http.Server
	DB               *sqlx.DB
	Consumer         *kafka.Subscriber
	ConsumerMessages <-chan *message.Message
	ConsumerHandler  *_internalKafka.ExerciseCompleteHandler
}

// NewApp создаёт новый экземпляр App
func NewApp(
	cfg *config.Config,
	httpServer *http.Server,
	db *sqlx.DB,
	consumer *kafka.Subscriber,
	consumerMessages <-chan *message.Message,
	consumerHandler *_internalKafka.ExerciseCompleteHandler,
) *App {
	return &App{
		Config:           cfg,
		HTTPServer:       httpServer,
		DB:               db,
		Consumer:         consumer,
		ConsumerMessages: consumerMessages,
		ConsumerHandler:  consumerHandler,
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

		// Kafka Consumer
		ProvideLogger,
		ProvideSubscriber,
		ProvideMessages,

		// Consumer Handler
		_internalKafka.NewExerciseCompleteHandler,

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

// ProvideLogger создает Watermill логгер
func ProvideLogger() watermill.LoggerAdapter {
	//return watermill.NewStdLogger(false, false)
	return watermill.NopLogger{}
}

// ProvideSubscriber создает Kafka Subscriber
func ProvideSubscriber(cfg *config.Config, logger watermill.LoggerAdapter) (*kafka.Subscriber, error) {
	return kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       []string{cfg.KafkaBroker},
			ConsumerGroup: cfg.KafkaExerciseCompletedGroup,
			Unmarshaler:   kafka.DefaultMarshaler{},
		},
		logger,
	)
}

// ProvideMessages создает канал сообщений
func ProvideMessages(cfg *config.Config, subscriber *kafka.Subscriber) (<-chan *message.Message, error) {
	ctx := context.Background()
	return subscriber.Subscribe(ctx, cfg.KafkaExerciseCompletedTopic)
}
