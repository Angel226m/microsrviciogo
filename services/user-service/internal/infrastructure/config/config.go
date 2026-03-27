package config

import "os"

type Config struct {
	Port           string
	DatabaseURL    string
	RedisURL       string
	NatsURL        string
	JaegerEndpoint string
	JWTSecret      string
	LogLevel       string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8081"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://cloudmart:cloudmart_secret@localhost:5432/cloudmart?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
		NatsURL:        getEnv("NATS_URL", "nats://localhost:4222"),
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:4318/v1/traces"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-key"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
