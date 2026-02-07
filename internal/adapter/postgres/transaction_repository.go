package postgres

import (
	"context"
	"time"
	model "transaction-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		"INSERT INTO transactions (id, name, amount, user_id, category_id, created_at = NOW()) VALUES ($1, $2, $3, $4, $5, $6)",
		transaction.ID, transaction.Name, transaction.Amount, transaction.UserID, transaction.CategoryID, transaction.CreatedAt,
	)
	return err
}

func (r *TransactionRepository) GetById(id uuid.UUID) (*model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction model.Transaction

	row := r.db.QueryRow(ctx,
		"SELECT id, name, amount, user_id, category_id, created_at FROM transactions WHERE id = $1",
		id,
	)

	err := row.Scan(&transaction.ID, &transaction.Name, &transaction.Amount, &transaction.UserID, &transaction.CategoryID, &transaction.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) GetAllByUserId(userId uuid.UUID) ([]model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transactions []model.Transaction

	rows, err := r.db.Query(ctx,
		`SELECT id, name, amount, user_id, category_id, created_at
         FROM transactions 
         WHERE user_id = $1
         ORDER BY created_at DESC`,
		userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Name,
			&transaction.Amount,
			&transaction.UserID,
			&transaction.CategoryID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(transaction *model.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`UPDATE transactions 
         SET name = $1, amount = $2, category_id = $3
         WHERE id = $4 AND user_id = $5`,
		transaction.Name,
		transaction.Amount,
		transaction.CategoryID,
		transaction.ID,
		transaction.UserID,
	)

	return err
}

func (r *TransactionRepository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		"DELETE FROM transactions WHERE id = $1",
		id,
	)

	return err
}
