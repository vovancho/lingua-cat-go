package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

type ServiceName string
type ExerciseCompletedTopic string

type Config struct {
	ServiceName                 ServiceName
	DBDSN                       string
	HTTPPort                    string
	AuthPublicKeyPath           string
	KafkaBroker                 string
	KafkaExerciseCompletedTopic string
	KafkaExerciseCompletedGroup string
	KeycloakAdminRealmEndpoint  string
	KeycloakAdminToken          string
	JaegerCollectorEndpoint     string
	Timeout                     int // in seconds
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
		ServiceName:                 ServiceName(os.Getenv("SERVICE_NAME")),
		DBDSN:                       os.Getenv("DB_DSN"),
		HTTPPort:                    os.Getenv("HTTP_PORT"),
		AuthPublicKeyPath:           os.Getenv("AUTH_PUBLIC_KEY_PATH"),
		KafkaBroker:                 os.Getenv("KAFKA_BROKER"),
		KafkaExerciseCompletedTopic: os.Getenv("KAFKA_EXERCISE_COMPLETED_TOPIC"),
		KafkaExerciseCompletedGroup: os.Getenv("KAFKA_EXERCISE_COMPLETED_GROUP"),
		KeycloakAdminRealmEndpoint:  os.Getenv("KEYCLOAK_ADMIN_REALM_ENDPOINT"),
		KeycloakAdminToken:          os.Getenv("KEYCLOAK_ADMIN_TOKEN"),
		JaegerCollectorEndpoint:     os.Getenv("JAEGER_COLLECTOR_ENDPOINT"),
		Timeout:                     timeout,
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}
