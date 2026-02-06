package main

import (
	"context"
	"log"
	"net/http"

	"transaction-service/internal/handler"
	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Репозитории
	transRepo := postgres.NewTransactionRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)

	// Сервис
	authSvc := service.NewAuthService(userRepo, tokenRepo, hasher, jwtSvc)

	// Хендлер
	authHandler := handler.NewAuthHandler(authSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.RegisterHandler)
	mux.HandleFunc("/auth/tokens", authHandler.AuthenticateHandler)
	mux.HandleFunc("/auth/refresh", authHandler.RefreshHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
