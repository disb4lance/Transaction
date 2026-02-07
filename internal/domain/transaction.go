package model

import (
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         uuid.UUID       `db:"id"`
	Name       string          `db:"name"`
	Amount     decimal.Decimal `db:"amount"`
	UserID     uuid.UUID       `db:"user_id"`
	CategoryID uuid.UUID       `db:"category_id"`
	CreatedAt  time.Time       `db:"created_at"`

	Category *Category `db:"-"`
}

type Category struct {
	Id          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}
