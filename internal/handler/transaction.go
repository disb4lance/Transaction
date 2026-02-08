package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"transaction-service/internal/service/dto"

	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, error)
	GetById(id uuid.UUID) (*dto.TransactionResponse, error)
	GetAllByUserId(userId uuid.UUID) ([]dto.TransactionResponse, error)
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

// GetById godoc
// @Summary Получить тразакцию по ID
// @Description Возвращает информацию о транзакции по её идентификатору
// @Tags transactions
// @Accept json
// @Produce json
// @Param id query string true "ID транзакции (UUID)"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} string "invalid ID format"
// @Failure 404 {object} string "category not found"
// @Failure 500 {object} string "internal server error"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetById(id)
	if err != nil {
		if err.Error() == "not found" || strings.Contains(err.Error(), "not found") {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetAll godoc
// @Summary Получить все транзакции пользователя
// @Description Возвращает список всех транзакций пользователя
// @Tags transactions
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя"
// @Success 200 {array} dto.TransactionResponse
// @Failure 500 {object} string "internal server error"
// @Router /transactions/user/{user_id} [get]
func (h *TransactionHandler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id") // TODO доставать id пользователя из токена
	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetAllByUserId(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
