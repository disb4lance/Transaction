package http

import (
	"net/http"
	"transaction-service/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(transactionHandler *handler.TransactionHandler, categoryHandler *handler.CategoryHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// versioning
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/categories", categoryHandler.Create)
		r.Get("/categories", categoryHandler.GetAll)
		r.Get("/categories/{id}", categoryHandler.GetById)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/transactions", transactionHandler.CreateTransaction)
		r.Get("/transactions/{id}", transactionHandler.GetById)
		r.Get("/transactions/user/{user_id}", transactionHandler.GetAllByUserId)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	return r
}
