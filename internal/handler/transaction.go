package handler

import (
	"encoding/json"
	"net/http"
	"transaction-service/internal/service/dto"
)

type TransactionService interface {
	CreateTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, error)
}

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(s TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

// CreateTransaction godoc
// @Summary Создание новой транзакции
// @Description Создает новую транзакцию
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body dto.TransactionRequest true "Данные транзакции"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} string "invalid body"
// @Router /transactions  [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateTransaction(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
