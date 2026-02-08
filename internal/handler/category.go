package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"transaction-service/internal/service/dto"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(req dto.CategoryRequest) (*dto.CategoryResponse, error)
	GetById(id uuid.UUID) (*dto.CategoryResponse, error)
	GetAll() ([]dto.CategoryResponse, error)
}

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(s CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// Create godoc
// @Summary Создание новой категории
// @Description Создает новую категорию для транзакций
// @Tags categories
// @Accept json
// @Produce json
// @Param request body dto.CategoryRequest true "Данные категории"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} string "invalid body"
// @Router /categories  [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateCategory(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetById godoc
// @Summary Получить категорию по ID
// @Description Возвращает информацию о категории по её идентификатору
// @Tags categories
// @Accept json
// @Produce json
// @Param id query string true "ID категории (UUID)"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} string "invalid ID format"
// @Failure 404 {object} string "category not found"
// @Failure 500 {object} string "internal server error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
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
// @Summary Получить все категории
// @Description Возвращает список всех доступных категорий
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {object} string "internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
