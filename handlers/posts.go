package handlers

import (
	"fmt"
	"strconv"

	"gofr-blog-service/models"
	"gofr-blog-service/services"

	"gofr.dev/pkg/gofr"
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
	if err := ph.validateUpdateRequest(req); err != nil {
		return ph.errorResponse("Validation failed", err), nil
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

// Private helper methods (decorator pattern for consistent behavior)

// parseCreateRequest parses create post request
func (ph *PostHandler) parseCreateRequest(ctx *gofr.Context, req *models.CreatePostRequest) error {
	if err := ctx.Bind(req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	return nil
}

// parseUpdateRequest parses update post request
func (ph *PostHandler) parseUpdateRequest(ctx *gofr.Context, req *models.UpdatePostRequest) error {
	if err := ctx.Bind(req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	return nil
}

// validateCreateRequest validates the create post request
func (ph *PostHandler) validateCreateRequest(req models.CreatePostRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) < 3 || len(req.Title) > 200 {
		return fmt.Errorf("title must be between 3 and 200 characters")
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	if len(req.Content) < 10 {
		return fmt.Errorf("content must be at least 10 characters")
	}
	if req.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	if req.AuthorID <= 0 {
		return fmt.Errorf("valid author ID is required")
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	validStatuses := []string{"draft", "published", "archived"}
	for _, status := range validStatuses {
		if req.Status == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s", req.Status)
}

// validateUpdateRequest validates the update post request
func (ph *PostHandler) validateUpdateRequest(req models.UpdatePostRequest) error {
	if req.Title != "" && (len(req.Title) < 3 || len(req.Title) > 200) {
		return fmt.Errorf("title must be between 3 and 200 characters")
	}
	if req.Content != "" && len(req.Content) < 10 {
		return fmt.Errorf("content must be at least 10 characters")
	}
	if req.Status != "" {
		validStatuses := []string{"draft", "published", "archived"}
		valid := false
		for _, status := range validStatuses {
			if req.Status == status {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid status: %s", req.Status)
		}
	}
	return nil
}

// extractIDParam extracts and validates ID parameter from URL
func (ph *PostHandler) extractIDParam(ctx *gofr.Context) (int, error) {
	idStr := ctx.PathParam("id")
	if idStr == "" {
		return 0, fmt.Errorf("missing post ID")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid post ID format: %s", idStr)
	}

	if id <= 0 {
		return 0, fmt.Errorf("post ID must be positive: %d", id)
	}

	return id, nil
}

// extractPaginationParams extracts pagination parameters with defaults
func (ph *PostHandler) extractPaginationParams(ctx *gofr.Context) (page, pageSize int) {
	page = 1
	pageSize = 10

	if pageStr := ctx.Param("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Param("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
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

// Middleware decorators can be added here for cross-cutting concerns:
// - Authentication
// - Rate limiting
// - Request logging
// - Input sanitization
// - Response compression
