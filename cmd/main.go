package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "transaction-service/docs"
	"transaction-service/internal/app"
	"transaction-service/internal/config"

	"github.com/joho/godotenv"
)

// @title Finance Tracker API
// @version 1.0
// @description API для управления личными финансами

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization
// @description Введите ваш JWT токен без префикса "Bearer"
func main() {
	envPath := filepath.Join("..", ".env")

	err := godotenv.Load(envPath)

	logger := log.New(os.Stdout, "[transaction-service] ", log.LstdFlags)

	cfg := config.Load()
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		// Изменено: cfg.Server.Env и cfg.Server.Port
		logger.Printf("starting server on :%s", cfg.APIPort)
		if err := application.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown:", err)
	}

	application.DB.Close()

	logger.Println("server exiting")
}
