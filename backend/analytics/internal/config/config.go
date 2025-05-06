package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	DBDSN                       string
	HTTPPort                    string
	AuthPublicKeyPath           string
	Timeout                     int // in seconds
	KafkaBroker                 string
	KafkaExerciseCompletedTopic string
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
		DBDSN:                       os.Getenv("DB_DSN"),
		HTTPPort:                    os.Getenv("HTTP_PORT"),
		AuthPublicKeyPath:           os.Getenv("AUTH_PUBLIC_KEY_PATH"),
		Timeout:                     timeout,
		KafkaBroker:                 os.Getenv("KAFKA_BROKER"),
		KafkaExerciseCompletedTopic: os.Getenv("KAFKA_EXERCISE_COMPLETED_TOPIC"),
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}
