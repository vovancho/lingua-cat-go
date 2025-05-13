//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"net/http"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_internalGrpc "github.com/vovancho/lingua-cat-go/dictionary/delivery/grpc"
	"github.com/vovancho/lingua-cat-go/dictionary/delivery/grpc/gen"
	_internalHttp "github.com/vovancho/lingua-cat-go/dictionary/delivery/http"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/config"
	_internalValidator "github.com/vovancho/lingua-cat-go/dictionary/internal/validator"
	"github.com/vovancho/lingua-cat-go/dictionary/repository/postgres"
	"github.com/vovancho/lingua-cat-go/dictionary/usecase"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"github.com/vovancho/lingua-cat-go/pkg/db"
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
	GRPCServer *grpc.Server
	DB         *sqlx.DB
	Tracer     *sdktrace.TracerProvider
}

// NewApp создаёт новый экземпляр App
func NewApp(
	cfg *config.Config,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	db *sqlx.DB,
	tracer *sdktrace.TracerProvider,
) *App {
	return &App{
		Config:     cfg,
		HTTPServer: httpServer,
		GRPCServer: grpcServer,
		DB:         db,
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

		// Репозиторий
		postgres.NewDictionaryRepository,

		// Use case
		usecase.NewDictionaryUseCase,

		// Tracing
		ProvideTracingServiceName,
		ProvideTracingEndpoint,
		tracing.NewTracer,

		// HTTP Delivery
		newHTTPServer,

		// gRPC Delivery
		newGRPCServer,

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

func ProvideTracingServiceName(cfg *config.Config) tracing.ServiceName {
	return tracing.ServiceName(cfg.ServiceName)
}

func ProvideTracingEndpoint(cfg *config.Config) tracing.Endpoint {
	return tracing.Endpoint(cfg.JaegerCollectorEndpoint)
}

// newHTTPServer создаёт новый HTTP-сервер
func newHTTPServer(
	cfg *config.Config,
	validator *validator.Validate,
	trans ut.Translator,
	authService *auth.AuthService,
	dictionaryUseCase domain.DictionaryUseCase,
) *http.Server {
	router := http.NewServeMux()
	_internalHttp.NewDictionaryHandler(router, dictionaryUseCase, validator)

	// Register gRPC-Gateway handlers
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := dictionary.RegisterDictionaryServiceHandlerFromEndpoint(context.Background(), gwmux, cfg.GRPCPort, opts); err != nil {
		panic(err)
	}

	mainMux := http.NewServeMux()
	mainMux.Handle("/grpc-gw-swagger.json", http.FileServer(http.Dir("docs")))
	mainMux.Handle("/grpc-gateway/", authService.AuthMiddleware(otelhttp.NewHandler(gwmux, "grpc-gateway")))

	mainMux.Handle("/swagger.json", http.FileServer(http.Dir("docs")))
	mainMux.Handle("/", response.ErrorMiddleware(authService.AuthMiddleware(otelhttp.NewHandler(router, "dictionary-http")), trans))

	return &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mainMux,
	}
}

// newGRPCServer создаёт новый gRPC-сервер
func newGRPCServer(
	authService *auth.AuthService,
	dictionaryUseCase domain.DictionaryUseCase,
) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			authService.AuthInterceptor,
		),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	dictionary.RegisterDictionaryServiceServer(grpcServer, _internalGrpc.NewDictionaryHandler(dictionaryUseCase))

	return grpcServer
}
