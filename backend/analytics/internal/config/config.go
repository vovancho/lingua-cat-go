package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ExerciseCompletedTopic string

type Config struct {
	ServiceName                 string
	DBDSN                       string
	HTTPPort                    string
	AuthPublicKeyPath           string
	KafkaBroker                 string
	KafkaExerciseCompletedTopic string
	KafkaExerciseCompletedGroup string
	KeycloakAdminTokenEndpoint  string
	KeycloakAdminRealmEndpoint  string
	KeycloakAdminClientID       string
	KeycloakAdminClientSecret   string
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
		AuthPublicKeyPath:           os.Getenv("AUTH_PUBLIC_KEY_PATH"),
		KafkaBroker:                 os.Getenv("KAFKA_BROKER"),
		KafkaExerciseCompletedTopic: os.Getenv("KAFKA_EXERCISE_COMPLETED_TOPIC"),
		KafkaExerciseCompletedGroup: os.Getenv("KAFKA_EXERCISE_COMPLETED_GROUP"),
		KeycloakAdminTokenEndpoint:  os.Getenv("KEYCLOAK_TOKEN_ENDPOINT"),
		KeycloakAdminRealmEndpoint:  os.Getenv("KEYCLOAK_ADMIN_REALM_ENDPOINT"),
		KeycloakAdminClientID:       os.Getenv("KEYCLOAK_ADMIN_CLIENT_ID"),
		KeycloakAdminClientSecret:   os.Getenv("KEYCLOAK_ADMIN_CLIENT_SECRET"),
		JaegerCollectorEndpoint:     os.Getenv("JAEGER_COLLECTOR_ENDPOINT"),
		Timeout:                     time.Duration(timeout) * time.Second,
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}
