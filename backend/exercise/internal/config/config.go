package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName                 string
	DBDSN                       string
	HTTPPort                    string
	DictionaryGRPCAddress       string
	AuthPublicKeyPath           string
	KafkaBroker                 string
	KafkaExerciseCompletedTopic string
	JaegerCollectorEndpoint     string
	Timeout                     time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Error loading .env file, using default config", "error", err)
	}

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("load config error %v", err)
	}

	cfg := &Config{
		ServiceName:                 os.Getenv("SERVICE_NAME"),
		DBDSN:                       os.Getenv("DB_DSN"),
		HTTPPort:                    os.Getenv("HTTP_PORT"),
		DictionaryGRPCAddress:       os.Getenv("DICTIONARY_GRPC_HOST"),
		AuthPublicKeyPath:           os.Getenv("AUTH_PUBLIC_KEY_PATH"),
		KafkaBroker:                 os.Getenv("KAFKA_BROKER"),
		KafkaExerciseCompletedTopic: os.Getenv("KAFKA_EXERCISE_COMPLETED_TOPIC"),
		JaegerCollectorEndpoint:     os.Getenv("JAEGER_COLLECTOR_ENDPOINT"),
		Timeout:                     time.Duration(timeout) * time.Second,
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}
