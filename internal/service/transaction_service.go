package service

import (
	model "transaction-service/internal/domain"
	"transaction-service/internal/service/dto"

	"github.com/google/uuid"
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

func (s *TransactionService) CreateTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, error) {
	transaction := &model.Transaction{
		ID:         uuid.New(),
		Name:       req.Name,
		Amount:     req.Amount,
		UserID:     uuid.New(), //TODO брать из токена
		CategoryID: req.CategoryID,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	return &dto.TransactionResponse{
		Id:         transaction.ID,
		Name:       transaction.Name,
		Amount:     req.Amount,
		UserID:     uuid.New(),
		CategoryID: req.CategoryID,
		CreatedAt:  transaction.CreatedAt,
	}, nil
}
