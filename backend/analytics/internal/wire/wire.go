//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"net/http"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	_internalHttp "github.com/vovancho/lingua-cat-go/analytics/delivery/http"
	_internalKafka "github.com/vovancho/lingua-cat-go/analytics/delivery/kafka"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/config"
	_internalValidator "github.com/vovancho/lingua-cat-go/analytics/internal/validator"
	"github.com/vovancho/lingua-cat-go/analytics/repository/clickhouse"
	_httpRepo "github.com/vovancho/lingua-cat-go/analytics/repository/http"
	usecase "github.com/vovancho/lingua-cat-go/analytics/usecase"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"github.com/vovancho/lingua-cat-go/pkg/db"
	"github.com/vovancho/lingua-cat-go/pkg/keycloak"
	"github.com/vovancho/lingua-cat-go/pkg/response"
	"github.com/vovancho/lingua-cat-go/pkg/tracing"
	"github.com/vovancho/lingua-cat-go/pkg/translator"
	_pkgValidator "github.com/vovancho/lingua-cat-go/pkg/validator"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// App представляет приложение с конфигурацией и серверами
type App struct {
	Config           *config.Config
	HTTPServer       *http.Server
	DB               *sqlx.DB
	Consumer         *kafka.Subscriber
	ConsumerMessages <-chan *message.Message
	ConsumerHandler  *_internalKafka.ExerciseCompleteHandler
	Tracer           *sdktrace.TracerProvider
}

// NewApp создаёт новый экземпляр App
func NewApp(
	cfg *config.Config,
	httpServer *http.Server,
	db *sqlx.DB,
	consumer *kafka.Subscriber,
	consumerMessages <-chan *message.Message,
	consumerHandler *_internalKafka.ExerciseCompleteHandler,
	tracer *sdktrace.TracerProvider,
) *App {
	return &App{
		Config:           cfg,
		HTTPServer:       httpServer,
		DB:               db,
		Consumer:         consumer,
		ConsumerMessages: consumerMessages,
		ConsumerHandler:  consumerHandler,
		Tracer:           tracer,
	}
}

func InitializeApp() (*App, error) {
	wire.Build(
		// Конфигурация
		config.Load,

		// База данных
		ProvideDriverName,
		ProvideDSN,
		db.NewDB,

		// Переводчик
		translator.NewTranslator,

		// Валидатор
		ProvideInternalValidator,

		// Аутентификация
		ProvidePublicKeyPath,
		auth.NewAuthService,

		// Репозиторий
		clickhouse.NewExerciseCompleteRepository,
		ProvideUserHttpClient,
		ProvideKeycloakAdminClient,
		_httpRepo.NewUserRepository,

		// Use case
		usecase.NewExerciseCompleteUseCase,
		usecase.NewUserUseCase,

		// Tracing
		ProvideTracingServiceName,
		ProvideTracingEndpoint,
		tracing.NewTracer,

		// Responder
		response.NewResponder,

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

func ProvideDriverName(cfg *config.Config) db.DriverName {
	return db.DriverName("clickhouse")
}

func ProvideDSN(cfg *config.Config) db.DSN {
	return db.DSN(cfg.DBDSN)
}

func ProvideInternalValidator(trans ut.Translator) *validator.Validate {
	pkgValidator, err := _pkgValidator.NewValidator(trans)
	if err != nil {
		panic(err)
	}

	internalValidator, err := _internalValidator.NewValidator(pkgValidator, trans)
	if err != nil {
		panic(err)
	}

	return internalValidator
}

func ProvidePublicKeyPath(cfg *config.Config) auth.PublicKeyPath {
	return auth.PublicKeyPath(cfg.AuthPublicKeyPath)
}

func ProvideTracingServiceName(cfg *config.Config) tracing.ServiceName {
	return tracing.ServiceName(cfg.ServiceName)
}

func ProvideTracingEndpoint(cfg *config.Config) tracing.Endpoint {
	return tracing.Endpoint(cfg.JaegerCollectorEndpoint)
}

// newHTTPServer создаёт новый HTTP-сервер
func newHTTPServer(
	cfg *config.Config,
	authService *auth.AuthService,
	exerciseCompleteUseCase domain.ExerciseCompleteUseCase,
	userUseCase domain.UserUseCase,
	responder response.Responder,
) *http.Server {
	router := http.NewServeMux()
	_internalHttp.NewExerciseCompleteHandler(router, responder, exerciseCompleteUseCase, userUseCase, authService)

	handler := authService.AuthMiddleware(router)
	handler = response.ErrorMiddleware(handler)
	handler = otelhttp.NewHandler(handler, "analytics-http")
	handler = http.TimeoutHandler(handler, cfg.Timeout, "Request timeout")

	mainMux := http.NewServeMux()
	mainMux.Handle("/swagger.json", http.FileServer(http.Dir("docs")))
	mainMux.Handle("/", handler)

	return &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mainMux,
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

// ProvideUserHttpClient создает HTTP-клиент для httpUserRepository
func ProvideUserHttpClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

// ProvideKeycloakAdminClient создает Keycloak AdminClient
func ProvideKeycloakAdminClient(cfg *config.Config, client *http.Client) *keycloak.AdminClient {
	return keycloak.NewAdminClient(
		keycloak.AdminClientConfig{
			TokenEndpoint:      cfg.KeycloakAdminTokenEndpoint,
			AdminRealmEndpoint: cfg.KeycloakAdminRealmEndpoint,
			ClientID:           cfg.KeycloakAdminClientID,
			ClientSecret:       cfg.KeycloakAdminClientSecret,
		},
		client,
	)
}
