package main

import (
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
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

	trans := initTranslator()
	validate := initValidator(trans)

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

func initTranslator() ut.Translator {
	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)
	trans, ok := uni.GetTranslator("ru")
	if !ok {
		panic("не удалось получить переводчик для ru")
	}
	return trans
}

func initValidator(trans ut.Translator) *validator.Validate {
	validate := validator.New()

	validate.RegisterTranslation("valid_dictionary_type", trans, func(ut ut.Translator) error {
		return ut.Add("valid_dictionary_type", "{0} должен быть валидным типом", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_dictionary_type", fe.Field())
		return t
	})
	if err := domain.RegisterDictionaryTypeValidation(validate); err != nil {
		panic(fmt.Errorf("failed to register dictionary type validation: %w", err))
	}
	validate.RegisterTranslation("valid_dictionary_lang", trans, func(ut ut.Translator) error {
		return ut.Add("valid_dictionary_lang", "{0} должен быть валидным языком", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_dictionary_lang", fe.Field())
		return t
	})
	if err := domain.RegisterDictionaryLangValidation(validate); err != nil {
		panic(fmt.Errorf("failed to register dictionary lang validation: %w", err))
	}
	if err := rutranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic("не удалось зарегистрировать русские переводы: " + err.Error())
	}
	return validate
}
