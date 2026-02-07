package dto

import "github.com/google/uuid"

type CategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type CategoryResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}
