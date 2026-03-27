// ═══════════════════════════════════════════════════════════════
// Configuration – Environment-driven config for API Gateway
// ═══════════════════════════════════════════════════════════════
package config

import "os"

type Config struct {
	Port                   string
	LogLevel               string
	RedisURL               string
	JaegerEndpoint         string
	JWTSecret              string
	UserServiceURL         string
	ProductServiceURL      string
	OrderServiceURL        string
	PaymentServiceURL      string
	InventoryServiceURL    string
	NotificationServiceURL string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		RedisURL:               getEnv("REDIS_URL", "localhost:6379"),
		JaegerEndpoint:         getEnv("JAEGER_ENDPOINT", "http://localhost:4318/v1/traces"),
		JWTSecret:              getEnv("JWT_SECRET", "dev-secret-key"),
		UserServiceURL:         getEnv("USER_SERVICE_URL", "http://localhost:8081"),
		ProductServiceURL:      getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082"),
		OrderServiceURL:        getEnv("ORDER_SERVICE_URL", "http://localhost:8083"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084"),
		InventoryServiceURL:    getEnv("INVENTORY_SERVICE_URL", "http://localhost:8085"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8086"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
