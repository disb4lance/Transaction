// config/config.go
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// App
	APIHost string
	APIPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	Redis RedisConfig

	// Kafka
	Kafka KafkaConfig

	// JWT
	JWTSecret            string
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers       string
	CreatedTopic  string
	UpdatedTopic  string
	DeletedTopic  string
	ConsumerGroup string
}

func Load() *Config {
	return &Config{
		// App
		APIHost: getEnv("API_HOST", "0.0.0.0"),
		APIPort: getEnv("API_PORT", "8081"),

		// DB
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "transaction"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "password"),
		DBName:     getEnv("POSTGRES_DB", "transaction_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Redis
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},

		// Kafka
		Kafka: KafkaConfig{
			Brokers:       getEnv("KAFKA_BROKERS", "localhost:9092"),
			CreatedTopic:  getEnv("KAFKA_CREATED_TOPIC", "transactions.created"),
			UpdatedTopic:  getEnv("KAFKA_UPDATED_TOPIC", "transactions.updated"),
			DeletedTopic:  getEnv("KAFKA_DELETED_TOPIC", "transactions.deleted"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "transaction-service"),
		},

		// JWT
		JWTSecret:            getEnv("JWT_SECRET", "secret"),
		JWTAccessExpiration:  getDuration("JWT_ACCESS_EXPIRATION", 15*time.Minute),
		JWTRefreshExpiration: getDuration("JWT_REFRESH_EXPIRATION", 24*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}
