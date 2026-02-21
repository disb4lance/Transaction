package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"transaction-service/internal/middleware"
	"transaction-service/internal/service/dto"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TransactionService interface {
	Create(req dto.TransactionRequest) (*dto.TransactionResponse, error)
	GetById(id uuid.UUID) (*dto.TransactionResponse, error)
	GetAllByUserId(userId uuid.UUID) ([]dto.TransactionResponse, error)
	Update(id uuid.UUID, userId uuid.UUID, transactionNew dto.EditTransactionRequest) error
	Delete(userId uuid.UUID, id uuid.UUID) error
}

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(s TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

// Create godoc
// @Summary Создание новой транзакции
// @Description Создает новую транзакцию
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body dto.TransactionRequest true "Данные транзакции"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} string "invalid body"
// @Security BearerAuth
// @Router /transactions  [post]
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Create(req)
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
// @Security BearerAuth
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetById(w http.ResponseWriter, r *http.Request) {

	_, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

// GetAllByUserId godoc
// @Summary Получить все транзакции пользователя
// @Description Возвращает список всех транзакций пользователя
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {array} dto.TransactionResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Security BearerAuth
// @Router /transactions [get]
func (h *TransactionHandler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := h.service.GetAllByUserId(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Update godoc
// @Summary Обновление транзакции
// @Description Обновляет существующую транзакцию по ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "ID транзакции" format(uuid)
// @Param request body dto.EditTransactionRequest true "Данные для обновления"
// @Success 204 "No Content - транзакция успешно обновлена"
// @Failure 400 {object} string "invalid request"
// @Failure 404 {object} string "transaction not found"
// @Failure 500 {object} string "internal server error"
// @Security BearerAuth
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.EditTransactionRequest

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	err = h.service.Update(transactionID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary Удаление транзакции
// @Description Удаляет существующую транзакцию по ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "ID транзакции (UUID)" format(uuid)
// @Success 204 "No Content - транзакция успешно удалена"
// @Failure 400 {object} string "invalid transaction id"
// @Failure 401 {object} string "unauthorized"
// @Failure 404 {object} string "transaction not found"
// @Failure 500 {object} string "internal server error"
// @Security BearerAuth
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(userID, transactionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
