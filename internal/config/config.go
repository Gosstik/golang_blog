package config

import "os"

type Config struct {
	GRPCPort    string
	HTTPPort    string
	PostgresDSN string
	RedisAddr   string
}

func Load() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", ":50051"),
		HTTPPort:    getEnv("HTTP_PORT", ":8090"),
		PostgresDSN: getEnv("POSTGRES_DSN", "host=localhost port=5432 user=blog password=blog dbname=blog sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
