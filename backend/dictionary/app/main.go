package main

import (
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"log/slog"
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

	validate := validator.New()
	trans := initTranslator(validate)

	router := http.NewServeMux()

	dictionaryRepo := _dictionaryRepo.NewPostgresDictionaryRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	du := _dictionaryUcase.NewDictionaryUseCase(dictionaryRepo, timeoutContext)
	_dictionaryHttpDelivery.NewDictionaryHandler(router, validate, du)

	server := http.Server{
		Addr:    ":80",
		Handler: response.ErrorMiddleware(router, trans),
	}
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

func initTranslator(v *validator.Validate) ut.Translator {
	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)
	trans, ok := uni.GetTranslator("ru")
	if !ok {
		panic("не удалось получить переводчик для ru")
	}
	if err := rutranslations.RegisterDefaultTranslations(v, trans); err != nil {
		panic("не удалось зарегистрировать русские переводы: " + err.Error())
	}
	return trans
}
