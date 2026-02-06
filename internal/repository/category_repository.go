package repository

import (
	"github.com/google/uuid"

	model "transaction-service/internal/models"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	GetById(id uuid.UUID) (*model.Category, error)
	GetAll() ([]model.Category, error)
}
