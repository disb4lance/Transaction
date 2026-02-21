package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"transaction-service/internal/adapter/kafka"
	"transaction-service/internal/adapter/postgres"
	"transaction-service/internal/adapter/redis"
	"transaction-service/internal/config"
	"transaction-service/internal/handler"
	"transaction-service/internal/service"
	transport "transaction-service/internal/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func New(cfg *config.Config, logger *log.Logger) (*App, error) {
	// Изменено: cfg.Database.* вместо cfg.DB.*
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	redisCache, err := redis.NewRedisClientWithError(cfg.Redis)
	if err != nil {
		return nil, err
	}

	producer := kafka.NewProducerWithBrokers(cfg.Kafka)

	// Репозитории
	transRepo := postgres.NewTransactionRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)

	// Сервис
	catSer := service.NewCategoryService(categoryRepo, redisCache)
	transSer := service.NewTransactionService(transRepo, redisCache, producer)

	// Хендлер
	categoryHandler := handler.NewCategoryHandler(catSer)
	transactionHandler := handler.NewTransactionHandler(transSer)

	// router
	router := transport.NewRouter(transactionHandler, categoryHandler, cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{
		Server: srv,
		DB:     db,
	}, nil
}
