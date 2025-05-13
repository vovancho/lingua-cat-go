//go:build wireinject
// +build wireinject

package wire

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	_internalHttp "github.com/vovancho/lingua-cat-go/exercise/delivery/http"
	_internalGrpc "github.com/vovancho/lingua-cat-go/exercise/repository/grpc"
	postgres "github.com/vovancho/lingua-cat-go/exercise/repository/postgres"
	usecase "github.com/vovancho/lingua-cat-go/exercise/usecase"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/config"
	_internalValidator "github.com/vovancho/lingua-cat-go/exercise/internal/validator"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"github.com/vovancho/lingua-cat-go/pkg/db"
	"github.com/vovancho/lingua-cat-go/pkg/eventpublisher"
	"github.com/vovancho/lingua-cat-go/pkg/response"
	"github.com/vovancho/lingua-cat-go/pkg/tracing"
	"github.com/vovancho/lingua-cat-go/pkg/translator"
	"github.com/vovancho/lingua-cat-go/pkg/txmanager"
	_pkgValidator "github.com/vovancho/lingua-cat-go/pkg/validator"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

// App представляет приложение с конфигурацией и серверами
type App struct {
	Config     *config.Config
	HTTPServer *http.Server
	DB         *sqlx.DB
	Outbox     *message.Router
	Tracer     *sdktrace.TracerProvider
}

// NewApp создаёт новый экземпляр App
func NewApp(
	cfg *config.Config,
	httpServer *http.Server,
	db *sqlx.DB,
	outbox *message.Router,
	tracer *sdktrace.TracerProvider,
) *App {
	return &App{
		Config:     cfg,
		HTTPServer: httpServer,
		DB:         db,
		Outbox:     outbox,
		Tracer:     tracer,
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
		txmanager.New,

		// Переводчик
		translator.NewTranslator,

		// Валидатор
		ProvideInternalValidator,

		// Аутентификация
		ProvidePublicKeyPath,
		auth.NewAuthService,

		// GRPC соединение
		ProvideGRPCConn,

		// Репозиторий
		postgres.NewExerciseRepository,
		postgres.NewTaskRepository,
		_internalGrpc.NewDictionaryRepository,

		// Use case
		ProvideExerciseCompletedPublisher,
		usecase.NewExerciseUseCase,
		usecase.NewTaskUseCase,
		usecase.NewDictionaryUseCase,

		// Tracing
		ProvideTracingServiceName,
		ProvideTracingEndpoint,
		tracing.NewTracer,

		// Responder
		response.NewResponder,

		// HTTP Delivery
		newHTTPServer,

		// Watermill Outbox
		ProvideLogger,
		ProvideSubscriber,
		ProvidePublisher,
		ProvideOutboxRouter,

		// App
		NewApp,
	)
	return &App{}, nil
}

func ProvideDriverName(cfg *config.Config) db.DriverName {
	return db.DriverName("postgres")
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

func ProvideExerciseCompletedPublisher(cfg *config.Config) usecase.ExerciseCompletedPublisherInterface {
	topic := eventpublisher.EventTopicName(cfg.KafkaExerciseCompletedTopic)

	return eventpublisher.NewEventPublisher(topic)
}

func ProvideTracingServiceName(cfg *config.Config) tracing.ServiceName {
	return tracing.ServiceName(cfg.ServiceName)
}

func ProvideTracingEndpoint(cfg *config.Config) tracing.Endpoint {
	return tracing.Endpoint(cfg.JaegerCollectorEndpoint)
}

func ProvideGRPCConn(cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cfg.DictionaryGRPCAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// newHTTPServer создаёт новый HTTP-сервер
func newHTTPServer(
	cfg *config.Config,
	validator *validator.Validate,
	authService *auth.AuthService,
	exerciseUseCase domain.ExerciseUseCase,
	taskUseCase domain.TaskUseCase,
	responder response.Responder,
) *http.Server {
	router := http.NewServeMux()
	_internalHttp.NewExerciseHandler(router, responder, exerciseUseCase, validator, authService)
	_internalHttp.NewTaskHandler(router, responder, taskUseCase, exerciseUseCase, validator, authService)

	mainMux := http.NewServeMux()
	mainMux.Handle("/swagger.json", http.FileServer(http.Dir("docs")))
	mainMux.Handle("/", otelhttp.NewHandler(response.ErrorMiddleware(authService.AuthMiddleware(router)), "exercise-http"))

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

// ProvideSubscriber создает SQL-подписчик
func ProvideSubscriber(db *sqlx.DB, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	return sql.NewSubscriber(
		db,
		sql.SubscriberConfig{
			SchemaAdapter:    sql.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   sql.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: false,
		},
		logger,
	)
}

// ProvidePublisher создает Kafka-паблишер
func ProvidePublisher(cfg *config.Config, logger watermill.LoggerAdapter) (message.Publisher, error) {
	return kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{cfg.KafkaBroker},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
}

// ProvideOutboxRouter создает outbox роутер из БД ProvideSubscriber в ProvidePublisher
func ProvideOutboxRouter(
	logger watermill.LoggerAdapter,
	subscriber message.Subscriber,
	publisher message.Publisher,
	cfg *config.Config,
) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddNoPublisherHandler(
		"outbox_to_kafka",
		cfg.KafkaExerciseCompletedTopic,
		subscriber,
		func(msg *message.Message) error {
			return publisher.Publish(cfg.KafkaExerciseCompletedTopic, msg)
		},
	)

	return router, nil
}
