package main

import (
	"fmt"
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

	slog.Warn("DB_DSN: ", "dbDsn", dbDsn)

	dbConn, err := sqlx.Open("postgres", dbDsn)
	if err != nil {
		slog.Error(err.Error())
	}
	err = dbConn.Ping()
	if err != nil {
		slog.Error(err.Error())
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	router := http.NewServeMux()

	dictionaryRepo := _dictionaryRepo.NewPostgresDictionaryRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	du := _dictionaryUcase.NewDictionaryUseCase(dictionaryRepo, timeoutContext)
	_dictionaryHttpDelivery.NewDictionaryHandler(router, du)

	server := http.Server{
		Addr:    ":80",
		Handler: router,
	}
	fmt.Println("Server is listening on port 80")
	server.ListenAndServe()
}
