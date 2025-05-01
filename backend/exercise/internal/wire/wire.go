//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	//"github.com/vovancho/lingua-cat-go/exercise/exercise/usecase"
	"github.com/vovancho/lingua-cat-go/exercise/exercise/grpc"
	"github.com/vovancho/lingua-cat-go/exercise/exercise/postgres"
	"github.com/vovancho/lingua-cat-go/exercise/internal/auth"
	"github.com/vovancho/lingua-cat-go/exercise/internal/config"
	"github.com/vovancho/lingua-cat-go/exercise/internal/db"
	"github.com/vovancho/lingua-cat-go/exercise/internal/translator"
	_internalValidator "github.com/vovancho/lingua-cat-go/exercise/internal/validator"
)

// App представляет приложение с конфигурацией и серверами
type App struct {
	Config *config.Config
	//HTTPServer *http.Server
	DB *sqlx.DB
}

// NewApp создаёт новый экземпляр App
func NewApp(cfg *config.Config /*, httpServer *http.Server*/, db *sqlx.DB) *App {
	return &App{
		Config: cfg,
		//HTTPServer: httpServer,
		DB: db,
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
		postgres.NewPostgresExerciseRepository,
		postgres.NewPostgresTaskRepository,
		grpc.NewPostgresDictionaryRepository,

		// Use case
		//ProvideUseCaseTimeout,
		//usecase.NewExerciseUseCase,

		// HTTP Delivery
		//newHTTPServer,

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
//func ProvideUseCaseTimeout(cfg *config.Config) usecase.Timeout {
//	return usecase.Timeout(time.Duration(cfg.Timeout) * time.Second)
//}

// newHTTPServer создаёт новый HTTP-сервер
//func newHTTPServer(
//	cfg *config.Config,
//	validate *validator.Validate,
//	trans ut.Translator,
//	authService *auth.AuthService,
//	exerciseUcase domain.ExerciseUseCase,
//) *http.Server {
//	router := http.NewServeMux()
//	_internalHttp.NewExerciseHandler(router, validate, exerciseUcase)
//	return &http.Server{
//		Addr:    cfg.HTTPPort,
//		Handler: response.ErrorMiddleware(authService.AuthMiddleware(router), trans),
//	}
//}
