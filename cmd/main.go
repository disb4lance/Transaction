package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
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
// @host localhost:8080
// @BasePath /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	logger := log.New(os.Stdout, "[transaction-service] ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		logger.Printf("starting %s server on %s", cfg.Env, application.Server.Addr)
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
