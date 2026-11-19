package model

import "time"

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SLAHours    int       `json:"sla_hours"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=255"`
	SLAHours    int    `json:"sla_hours" validate:"required,min=1,max=720"`
}

type ReplaceCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=255"`
	SLAHours    int    `json:"sla_hours" validate:"required,min=1,max=720"`
	IsActive    bool   `json:"is_active"`
}

type PatchCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
	SLAHours    *int    `json:"sla_hours,omitempty" validate:"omitempty,min=1,max=720"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
