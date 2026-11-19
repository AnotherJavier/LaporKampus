package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,username"`
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,strongpassword"`
	Role     string `json:"role" validate:"omitempty,max=20"`
}

type ReplaceUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,username"`
	Email    string `json:"email" validate:"required,email,max=120"`
	Role     string `json:"role" validate:"required,max=20"`
	IsActive bool   `json:"is_active"`
}

type PatchUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=30,username"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email,max=120"`
	Role     *string `json:"role,omitempty" validate:"omitempty,max=20"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type AssignRoleRequest struct {
	Role string `json:"role" validate:"required,max=20"`
}
