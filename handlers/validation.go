package handlers

import (
	"errors"
	"strconv"
	"regexp"

	"gofr-blog-service/models"

	"gofr.dev/pkg/gofr"
)

// parseCreateRequest parses create post request
func (ph *PostHandler) parseCreateRequest(ctx *gofr.Context, req *models.CreatePostRequest) error {
	if err := ctx.Bind(req); err != nil {
		return errors.Join(errInvalidRequest, err)
	}
	return nil
}

// parseUpdateRequest parses update post request
func (ph *PostHandler) parseUpdateRequest(ctx *gofr.Context, req *models.UpdatePostRequest) error {
	if err := ctx.Bind(req); err != nil {
		return errors.Join(errInvalidRequest, err)
	}
	return nil
}

// validateCreateRequest validates the create post request
func (ph *PostHandler) validateCreateRequest(req models.CreatePostRequest) error {
	if req.Title == "" {
		return errors.Join(errValidation, errors.New("title is required"))
	}
	if len(req.Title) < 3 || len(req.Title) > 200 {
		return errors.Join(errValidation, errors.New("title must be between 3 and 200 characters"))
	}
	if req.Content == "" {
		return errors.Join(errValidation, errors.New("content is required"))
	}
	if len(req.Content) < 10 {
		return errors.Join(errValidation, errors.New("content must be at least 10 characters"))
	}
	if req.Slug == "" {
		return errors.Join(errValidation, errors.New("slug is required"))
	}
	if req.AuthorID <= 0 {
		return errors.Join(errValidation, errors.New("valid author ID is required"))
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
	return errors.Join(errValidation, errors.New("invalid status: "+req.Status))
}

// validateUpdateRequest validates the update post request
func (ph *PostHandler) validateUpdateRequest(req models.UpdatePostRequest) error {
	if req.Title != "" && (len(req.Title) < 3 || len(req.Title) > 200) {
		return errors.Join(errValidation, errors.New("title must be between 3 and 200 characters"))
	}
	if req.Content != "" && len(req.Content) < 10 {
		return errors.Join(errValidation, errors.New("content must be at least 10 characters"))
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
			return errors.Join(errValidation, errors.New("invalid status: "+req.Status))
		}
	}
	return nil
}

// extractIDParam extracts and validates ID parameter from URL
func (ph *PostHandler) extractIDParam(ctx *gofr.Context) (int, error) {
	idStr := ctx.PathParam("id")
	if idStr == "" {
		return 0, errors.Join(errInvalidID, errors.New("missing post ID"))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.Join(errInvalidID, errors.New("invalid post ID format: "+idStr))
	}

	if id <= 0 {
		return 0, errors.Join(errInvalidID, errors.New("post ID must be positive"))
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

// ValidateRegisterRequest validates the register request
func ValidateRegisterRequest(req *models.RegisterRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !isValidUsername(req.Username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if len(req.FirstName) < 2 || len(req.FirstName) > 50 {
		return errors.New("first name must be between 2 and 50 characters")
	}

	if req.LastName == "" {
		return errors.New("last name is required")
	}
	if len(req.LastName) < 2 || len(req.LastName) > 50 {
		return errors.New("last name must be between 2 and 50 characters")
	}

	if req.Bio != "" && len(req.Bio) > 500 {
		return errors.New("bio cannot exceed 500 characters")
	}

	if req.AvatarURL != "" && !isValidURL(req.AvatarURL) {
		return errors.New("invalid avatar URL format")
	}

	return nil
}

// ValidateLoginRequest validates the login request
func ValidateLoginRequest(req *models.LoginRequest) error {
	if req.Username == "" {
		return errors.New("username or email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// ValidateUpdateAuthorRequest validates the update author request
func ValidateUpdateAuthorRequest(req *models.UpdateAuthorRequest) error {
	if req.Email != "" && !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if req.FirstName != "" && (len(req.FirstName) < 2 || len(req.FirstName) > 50) {
		return errors.New("first name must be between 2 and 50 characters")
	}

	if req.LastName != "" && (len(req.LastName) < 2 || len(req.LastName) > 50) {
		return errors.New("last name must be between 2 and 50 characters")
	}

	if req.Bio != nil && len(*req.Bio) > 500 {
		return errors.New("bio cannot exceed 500 characters")
	}

	if req.AvatarURL != nil && *req.AvatarURL != "" && !isValidURL(*req.AvatarURL) {
		return errors.New("invalid avatar URL format")
	}

	return nil
}

// Helper validation functions
func isValidUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

func isValidEmail(email string) bool {
	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidURL(url string) bool {
	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}
