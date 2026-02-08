package service

import (
	"database/sql"
	"time"
	model "transaction-service/internal/domain"
	"transaction-service/internal/service/dto"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	GetById(id uuid.UUID) (*model.Category, error)
	GetAll() ([]model.Category, error)
}

type RedisClient interface {
	Get(key string, dest interface{}) bool
	Set(key string, value interface{}, ttl time.Duration) bool
	Delete(key string) bool
}

type CategoryService struct {
	categoryRepo CategoryRepository
	redis        RedisClient
}

func NewCategoryService(
	c CategoryRepository,
	r RedisClient,
) *CategoryService {
	return &CategoryService{
		categoryRepo: c,
		redis:        r,
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

	cacheKey := "categories:all"

	if s.redis != nil {
		var cachedResult []dto.CategoryResponse
		if s.redis.Get(cacheKey, &cachedResult) {
			return cachedResult, nil
		}
	}

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

	if s.redis != nil {
		s.redis.Set(cacheKey, result, 5*time.Minute)
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
