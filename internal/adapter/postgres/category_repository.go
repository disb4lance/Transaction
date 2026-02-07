package postgres

import (
	"context"
	"time"
	model "transaction-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *model.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`INSERT INTO categories (id, name, description)
		VALUES ($1, $2, $3)`,
		category.Id, category.Name, category.Description,
	)
	return err
}

func (r *CategoryRepository) GetById(id uuid.UUID) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var category model.Category
	row := r.db.QueryRow(ctx,
		`SELECT id, name, description
		 FROM categories 
		 WHERE id = $1`,
		id,
	)

	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var categories []model.Category

	rows, err := r.db.Query(ctx,
		`SELECT id, name, description 
         FROM categories 
         ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
