package kmodel

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionEvent struct {
	TransactionID string          `json:"transaction_id"`
	UserID        string          `json:"user_id"`
	CategoryID    string          `json:"category_id"`
	Amount        decimal.Decimal `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
}
