package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServiceName             string
	DBDSN                   string
	HTTPPort                string
	GRPCPort                string
	AuthPublicKeyPath       string
	JaegerCollectorEndpoint string
	Timeout                 int // in seconds
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
		ServiceName:             os.Getenv("SERVICE_NAME"),
		DBDSN:                   os.Getenv("DB_DSN"),
		HTTPPort:                os.Getenv("HTTP_PORT"),
		GRPCPort:                os.Getenv("GRPC_PORT"),
		AuthPublicKeyPath:       os.Getenv("AUTH_PUBLIC_KEY_PATH"),
		JaegerCollectorEndpoint: os.Getenv("JAEGER_COLLECTOR_ENDPOINT"),
		Timeout:                 timeout,
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}
