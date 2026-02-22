package http

import (
	"net/http"
	"transaction-service/internal/config"
	"transaction-service/internal/handler"
	mid "transaction-service/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	transactionHandler *handler.TransactionHandler,
	categoryHandler *handler.CategoryHandler,
	cfg *config.Config,
) *chi.Mux {

	authMiddleware := mid.NewAuthMiddleware(cfg.JWTSecret)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/categories", categoryHandler.GetAll)
		r.Get("/categories/{id}", categoryHandler.GetById)
		r.Post("/categories", categoryHandler.Create)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Post("/transactions", transactionHandler.Create)
			r.Get("/transactions/{id}", transactionHandler.GetById)
			r.Get("/transactions", transactionHandler.GetAllByUserId)
			r.Put("/transactions/{id}", transactionHandler.Update)
			r.Delete("/transactions/{id}", transactionHandler.Delete)
		})
	})

	return r
}
