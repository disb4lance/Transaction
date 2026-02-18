package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	ctx := context.Background()
	db, err := pgxpool.New(ctx, "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisCache := redis.NewRedisClient(redisAddr)

	producer := kafka.NewProducer(
		"localhost:9092",
		"transactions.created",
	)

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
	router := transport.NewRouter(transactionHandler, categoryHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &App{
		Server: srv,
		DB:     db,
	}, nil
}
