package handlers

import (
	"errors"

	"gofr-blog-service/models"
	"gofr-blog-service/services"

	"gofr.dev/pkg/gofr"
)

// Error definitions
var (
	ErrInvalidRequest = errors.New("invalid request format")
	ErrValidation     = errors.New("validation failed")
	ErrInvalidID      = errors.New("invalid post ID")
	ErrNotFound       = errors.New("post not found")
)

// PostHandler handles HTTP requests for posts with decorators pattern
type PostHandler struct {
	postService *services.PostService
}

// NewPostHandler creates a new post handler instance (dependency injection decorator)
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost handles POST /posts (HTTP decorator pattern)
func (ph *PostHandler) CreatePost(ctx *gofr.Context) (any, error) {
	// Request parsing decorator
	var req models.CreatePostRequest
	if err := ph.parseCreateRequest(ctx, &req); err != nil {
		return ph.errorResponse("Invalid request format", err), nil
	}

	// Validation decorator - moved from service to handler
	if err := ph.validateCreateRequest(req); err != nil {
		return ph.errorResponse("Validation failed", err), nil
	}

	// Business logic delegation decorator
	post, err := ph.postService.CreatePost(ctx, req)
	if err != nil {
		return ph.errorResponse("Failed to create post", err), nil
	}

	// Success response decorator
	return ph.successResponse("Post created successfully", post), nil
}

// GetPost handles GET /posts/{id}
func (ph *PostHandler) GetPost(ctx *gofr.Context) (any, error) {
	// Parameter extraction decorator
	id, err := ph.extractIDParam(ctx)
	if err != nil {
		return ph.errorResponse("Invalid post ID", err), nil
	}

	// Service call decorator
	post, err := ph.postService.GetPost(ctx, id)
	if err != nil {
		return ph.errorResponse("Post not found", err), nil
	}

	return ph.successResponse("Post retrieved successfully", post), nil
}

// ListPosts handles GET /posts with pagination
func (ph *PostHandler) ListPosts(ctx *gofr.Context) (any, error) {
	// Query parameter extraction decorator
	page, pageSize := ph.extractPaginationParams(ctx)

	// Service call with error handling decorator
	posts, err := ph.postService.ListPosts(ctx, page, pageSize)
	if err != nil {
		return ph.errorResponse("Failed to retrieve posts", err), nil
	}

	return ph.successResponse("Posts retrieved successfully", posts), nil
}

// UpdatePost handles PUT /posts/{id}
func (ph *PostHandler) UpdatePost(ctx *gofr.Context) (any, error) {
	// Parameter extraction decorator
	id, err := ph.extractIDParam(ctx)
	if err != nil {
		return ph.errorResponse("Invalid post ID", err), nil
	}

	// Request parsing decorator
	var req models.UpdatePostRequest
	if parseErr := ph.parseUpdateRequest(ctx, &req); parseErr != nil {
		return ph.errorResponse("Invalid request format", parseErr), nil
	}

	// Validation decorator - moved from service to handler
	if validateErr := ph.validateUpdateRequest(req); validateErr != nil {
		return ph.errorResponse("Validation failed", validateErr), nil
	}

	// Service call decorator
	post, err := ph.postService.UpdatePost(ctx, id, req)
	if err != nil {
		return ph.errorResponse("Failed to update post", err), nil
	}

	return ph.successResponse("Post updated successfully", post), nil
}

// DeletePost handles DELETE /posts/{id}
func (ph *PostHandler) DeletePost(ctx *gofr.Context) (any, error) {
	// Parameter extraction decorator
	id, err := ph.extractIDParam(ctx)
	if err != nil {
		return ph.errorResponse("Invalid post ID", err), nil
	}

	// Service call decorator
	err = ph.postService.DeletePost(ctx, id)
	if err != nil {
		return ph.errorResponse("Failed to delete post", err), nil
	}

	return ph.successResponse("Post deleted successfully", map[string]any{
		"deleted_id": id,
	}), nil
}

// Response formatting decorators for consistent API responses

// successResponse creates a standardized success response
func (ph *PostHandler) successResponse(message string, data any) map[string]any {
	return map[string]any{
		"success": true,
		"message": message,
		"data":    data,
	}
}

// errorResponse creates a standardized error response
func (ph *PostHandler) errorResponse(message string, err error) map[string]any {
	response := map[string]any{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	return response
}
