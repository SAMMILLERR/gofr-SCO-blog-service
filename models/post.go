package models

import (
	"time"
)

// Post represents a blog post in the system
type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title" validate:"required,min=3,max=200"`
	Content   string    `json:"content" db:"content" validate:"required,min=10"`
	Slug      string    `json:"slug" db:"slug" validate:"required,min=3,max=200"`
	AuthorID  int       `json:"author_id" db:"author_id" validate:"required"`
	Status    string    `json:"status" db:"status" validate:"required,oneof=draft published archived"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Title    string `json:"title" validate:"required,min=3,max=200"`
	Content  string `json:"content" validate:"required,min=10"`
	Slug     string `json:"slug" validate:"required,min=3,max=200"`
	AuthorID int    `json:"author_id" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=draft published archived"`
}

// UpdatePostRequest represents the request body for updating a post
type UpdatePostRequest struct {
	Title   string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Content string `json:"content,omitempty" validate:"omitempty,min=10"`
	Slug    string `json:"slug,omitempty" validate:"omitempty,min=3,max=200"`
	Status  string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

// PostListResponse represents the response for listing posts
type PostListResponse struct {
	Posts      []Post `json:"posts"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}
