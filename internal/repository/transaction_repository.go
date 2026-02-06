package repository

import (
	"github.com/google/uuid"

	model "transaction-service/internal/models"
)

type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetById(id uuid.UUID) (*model.Transaction, error)
	GetAllByUserId(userId uuid.UUID) ([]model.Transaction, error)
	Update(transaction *model.Transaction) error
	Delete(id uuid.UUID) error
}
