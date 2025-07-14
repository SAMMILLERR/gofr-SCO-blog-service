package models

import (
	"time"
)

// Author represents an author in the system
type Author struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Internal field, not exposed in JSON
	FirstName    string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName     string    `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Bio          string    `json:"bio" db:"bio" validate:"max=500"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url" validate:"url"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAuthorRequest represents the request body for creating an author
type CreateAuthorRequest struct {
	Username     string `json:"username" validate:"required,min=3,max=50"`
	Email        string `json:"email" validate:"required,email"`
	PasswordHash string `json:"-"` // Internal field, not exposed in JSON
	FirstName    string `json:"first_name" validate:"required,min=2,max=50"`
	LastName     string `json:"last_name" validate:"required,min=2,max=50"`
	Bio          string `json:"bio,omitempty" validate:"max=500"`
	AvatarURL    string `json:"avatar_url,omitempty" validate:"url"`
}

// UpdateAuthorRequest represents the request body for updating an author
type UpdateAuthorRequest struct {
	Email     string  `json:"email,omitempty" validate:"omitempty,email"`
	FirstName string  `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  string  `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Bio       *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

// AuthorResponse represents the public author response (without sensitive fields)
type AuthorResponse struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Bio        string    `json:"bio"`
	AvatarURL  string    `json:"avatar_url"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"` // Can be username or email
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Bio       string `json:"bio,omitempty" validate:"max=500"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"url"`
}

// LoginResponse represents the authentication response
type LoginResponse struct {
	Token  string         `json:"token"`
	Author AuthorResponse `json:"author"`
}
