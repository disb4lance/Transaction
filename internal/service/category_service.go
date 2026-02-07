package service

import (
	"github.com/google/uuid"

	"database/sql"
	model "transaction-service/internal/domain"
	"transaction-service/internal/service/dto"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	GetById(id uuid.UUID) (*model.Category, error)
	GetAll() ([]model.Category, error)
}

type CategoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(
	c CategoryRepository,
) *CategoryService {
	return &CategoryService{
		categoryRepo: c,
	}
}

func (s *CategoryService) CreateCategory(req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	category := &model.Category{
		Id:   uuid.New(),
		Name: req.Name,
	}

	if req.Description != nil && *req.Description != "" {
		category.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Id:          category.Id,
		Name:        category.Name,
		Description: req.Description,
	}, nil
}

func (s *CategoryService) GetAll() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.CategoryResponse, len(categories))

	for i, category := range categories {
		result[i] = dto.CategoryResponse{
			Id:          category.Id,
			Name:        category.Name,
			Description: nullStringToPtr(category.Description),
		}
	}

	return result, nil
}
func (s *CategoryService) GetById(id uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Id:          category.Id,
		Name:        category.Name,
		Description: nullStringToPtr(category.Description),
	}, nil
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
