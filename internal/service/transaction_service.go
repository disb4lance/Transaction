package service

import (
	"github.com/google/uuid"

	model "transaction-service/internal/domain"
)

type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetById(id uuid.UUID) (*model.Transaction, error)
	GetAllByUserId(userId uuid.UUID) ([]model.Transaction, error)
	Update(transaction *model.Transaction) error
	Delete(id uuid.UUID) error
}

type TransactionService struct {
	transactionRepo TransactionRepository
}

func NewTransactionService(
	t TransactionRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: t,
	}
}
