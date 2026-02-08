package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	Name       string          `json:"name" validate:"required,min=2,max=20"`
	Amount     decimal.Decimal `json:"amount"`
	CategoryID uuid.UUID       `json:"category_id"`
}

type TransactionResponse struct {
	Id         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Amount     decimal.Decimal `json:"amount"`
	UserID     uuid.UUID       `json:"user_id"`
	CategoryID uuid.UUID       `json:"category_id"`
	CreatedAt  time.Time       `json:"created_at"`
}
