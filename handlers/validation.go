package handlers

import (
	"errors"
	"strconv"

	"gofr-blog-service/models"

	"gofr.dev/pkg/gofr"
)

// parseCreateRequest parses create post request
func (ph *PostHandler) parseCreateRequest(ctx *gofr.Context, req *models.CreatePostRequest) error {
	if err := ctx.Bind(req); err != nil {
		return errors.Join(ErrInvalidRequest, err)
	}
	return nil
}

// parseUpdateRequest parses update post request
func (ph *PostHandler) parseUpdateRequest(ctx *gofr.Context, req *models.UpdatePostRequest) error {
	if err := ctx.Bind(req); err != nil {
		return errors.Join(ErrInvalidRequest, err)
	}
	return nil
}

// validateCreateRequest validates the create post request
func (ph *PostHandler) validateCreateRequest(req models.CreatePostRequest) error {
	if req.Title == "" {
		return errors.Join(ErrValidation, errors.New("title is required"))
	}
	if len(req.Title) < 3 || len(req.Title) > 200 {
		return errors.Join(ErrValidation, errors.New("title must be between 3 and 200 characters"))
	}
	if req.Content == "" {
		return errors.Join(ErrValidation, errors.New("content is required"))
	}
	if len(req.Content) < 10 {
		return errors.Join(ErrValidation, errors.New("content must be at least 10 characters"))
	}
	if req.Slug == "" {
		return errors.Join(ErrValidation, errors.New("slug is required"))
	}
	if req.AuthorID <= 0 {
		return errors.Join(ErrValidation, errors.New("valid author ID is required"))
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
	return errors.Join(ErrValidation, errors.New("invalid status: "+req.Status))
}

// validateUpdateRequest validates the update post request
func (ph *PostHandler) validateUpdateRequest(req models.UpdatePostRequest) error {
	if req.Title != "" && (len(req.Title) < 3 || len(req.Title) > 200) {
		return errors.Join(ErrValidation, errors.New("title must be between 3 and 200 characters"))
	}
	if req.Content != "" && len(req.Content) < 10 {
		return errors.Join(ErrValidation, errors.New("content must be at least 10 characters"))
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
			return errors.Join(ErrValidation, errors.New("invalid status: "+req.Status))
		}
	}
	return nil
}

// extractIDParam extracts and validates ID parameter from URL
func (ph *PostHandler) extractIDParam(ctx *gofr.Context) (int, error) {
	idStr := ctx.PathParam("id")
	if idStr == "" {
		return 0, errors.Join(ErrInvalidID, errors.New("missing post ID"))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.Join(ErrInvalidID, errors.New("invalid post ID format: "+idStr))
	}

	if id <= 0 {
		return 0, errors.Join(ErrInvalidID, errors.New("post ID must be positive"))
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
