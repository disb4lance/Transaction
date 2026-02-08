package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "transaction-service/docs"
	"transaction-service/internal/adapter/postgres"
	"transaction-service/internal/adapter/redis"
	"transaction-service/internal/handler"
	"transaction-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Finance Tracker API
// @version 1.0
// @description API для управления личными финансами
// @host localhost:8080
// @BasePath /api/v1
func main() {
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

	// Репозитории
	transRepo := postgres.NewTransactionRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)

	// Сервис
	catSer := service.NewCategoryService(categoryRepo, redisCache)
	transSer := service.NewTransactionService(transRepo)

	// Хендлер
	categoryHandler := handler.NewCategoryHandler(catSer)
	transactionHandler := handler.NewTransactionHandler(transSer)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/categories", categoryHandler.Create)
		r.Get("/categories", categoryHandler.GetAll)
		r.Get("/categories/{id}", categoryHandler.GetById)
		r.Post("/transactions", transactionHandler.CreateTransaction)
		r.Get("/transactions/{id}", transactionHandler.GetById)
		r.Get("/transactions/user/{user_id}", transactionHandler.GetAllByUserId)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
