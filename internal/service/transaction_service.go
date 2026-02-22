package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"transaction-service/internal/adapter/kafka/kmodel"
	model "transaction-service/internal/domain"
	"transaction-service/internal/service/contract"
	"transaction-service/internal/service/dto"

	"github.com/google/uuid"
)

type TransactionEventPublisher interface {
	PublishTransactionEvent(
		ctx context.Context,
		topic string,
		event kmodel.TransactionEvent,
	) error
}

type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetById(id uuid.UUID) (*model.Transaction, error)
	GetAllByUserId(userId uuid.UUID) ([]model.Transaction, error)
	Update(transaction *model.Transaction) error
	Delete(id uuid.UUID) error
}

type TransactionService struct {
	transactionRepo TransactionRepository
	redis           contract.RedisClient
	eventPublisher  TransactionEventPublisher
}

func NewTransactionService(
	t TransactionRepository,
	r contract.RedisClient,
	publisher TransactionEventPublisher,
) *TransactionService {
	return &TransactionService{
		transactionRepo: t,
		redis:           r,
		eventPublisher:  publisher,
	}
}

func (s *TransactionService) Create(userId uuid.UUID,
	req dto.TransactionRequest,
) (*dto.TransactionResponse, error) {

	transaction := &model.Transaction{
		ID:         uuid.New(),
		Name:       req.Name,
		Amount:     req.Amount,
		UserID:     userId,
		CategoryID: req.CategoryID,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	event := kmodel.TransactionEvent{
		TransactionID: transaction.ID.String(),
		UserID:        transaction.UserID.String(),
		CategoryID:    transaction.CategoryID.String(),
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
	}

	if err := s.eventPublisher.PublishTransactionEvent(
		context.Background(),
		"transactions.created",
		event,
	); err != nil {
		return nil, err
	}

	return &dto.TransactionResponse{
		Id:         transaction.ID,
		Name:       transaction.Name,
		Amount:     transaction.Amount,
		UserID:     transaction.UserID,
		CategoryID: transaction.CategoryID,
		CreatedAt:  transaction.CreatedAt,
	}, nil
}

func (s *TransactionService) GetById(id uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, fmt.Errorf("transaction with id %s not found", id)
	}

	return &dto.TransactionResponse{
		Id:         transaction.ID,
		Name:       transaction.Name,
		Amount:     transaction.Amount,
		UserID:     transaction.UserID,
		CategoryID: transaction.CategoryID,
		CreatedAt:  transaction.CreatedAt,
	}, nil
}

func (s *TransactionService) GetAllByUserId(userId uuid.UUID) ([]dto.TransactionResponse, error) {
	cacheKey := "transactions:all"

	if s.redis != nil {
		var cachedResult []dto.TransactionResponse
		if s.redis.Get(cacheKey, &cachedResult) {
			return cachedResult, nil
		}
	}

	transactions, err := s.transactionRepo.GetAllByUserId(userId)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TransactionResponse, len(transactions))

	for i, transaction := range transactions {
		result[i] = dto.TransactionResponse{
			Id:         transaction.ID,
			Name:       transaction.Name,
			Amount:     transaction.Amount,
			UserID:     transaction.UserID,
			CategoryID: transaction.CategoryID,
			CreatedAt:  transaction.CreatedAt,
		}
	}

	if s.redis != nil {
		s.redis.Set(cacheKey, result, 5*time.Minute)
	}

	return result, nil
}

func (s *TransactionService) Update(id uuid.UUID, userId uuid.UUID, transactionNew dto.EditTransactionRequest) error {
	transaction, err := s.transactionRepo.GetById(id)
	if err != nil {
		return err
	}

	if transaction == nil {
		return errors.New("transaction not found")
	}

	if transaction.UserID != userId {
		return errors.New("access denied")
	}

	transaction = &model.Transaction{
		ID:         id,
		Name:       transactionNew.Name,
		Amount:     transactionNew.Amount,
		UserID:     userId, // TODO из токена
		CategoryID: transactionNew.CategoryID,
	}

	event := kmodel.TransactionEvent{
		TransactionID: transaction.ID.String(),
		UserID:        transaction.UserID.String(),
		CategoryID:    transaction.CategoryID.String(),
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
	}

	err = s.transactionRepo.Update(transaction)
	if err != nil {
		return err
	}

	if err := s.eventPublisher.PublishTransactionEvent(
		context.Background(),
		"transactions.updated",
		event,
	); err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) Delete(userId uuid.UUID, id uuid.UUID) error {
	transaction, err := s.transactionRepo.GetById(id)
	if err != nil {
		return err
	}

	if transaction == nil {
		return errors.New("transaction not found")
	}

	if transaction.UserID != userId {
		return errors.New("access denied")
	}

	err = s.transactionRepo.Delete(id)
	if err != nil {
		return err
	}

	event := kmodel.TransactionEvent{
		TransactionID: transaction.ID.String(),
		UserID:        transaction.UserID.String(),
		CategoryID:    transaction.CategoryID.String(),
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
	}

	if err := s.eventPublisher.PublishTransactionEvent(
		context.Background(),
		"transactions.deleted",
		event,
	); err != nil {
		return err
	}

	return nil
}
