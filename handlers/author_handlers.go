package handlers

import (
	"errors"
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr-blog-service/interfaces"
	"gofr-blog-service/models"
)

type AuthorHandler struct {
	service interfaces.AuthorServiceInterface
}

func NewAuthorHandler(service interfaces.AuthorServiceInterface) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

func (h *AuthorHandler) Register(ctx *gofr.Context) (interface{}, error) {
	var req models.RegisterRequest

	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse("Invalid request format", err), nil
	}

	if err := ValidateRegisterRequest(&req); err != nil {
		return h.errorResponse("Validation failed", err), nil
	}

	author, err := h.service.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, errors.New("username already exists")) ||
			errors.Is(err, errors.New("email already exists")) {
			return h.errorResponse("Username or email already exists", err), nil
		}
		if errors.Is(err, errors.New("password must be at least 8 characters long")) {
			return h.errorResponse("Password validation failed", err), nil
		}
		return h.errorResponse("Failed to create account", err), nil
	}

	return map[string]interface{}{
		"success": true,
		"message": "Account created successfully",
		"author":  author,
	}, nil
}

func (h *AuthorHandler) Login(ctx *gofr.Context) (interface{}, error) {
	var req models.LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse("Invalid request format", err), nil
	}

	if err := ValidateLoginRequest(&req); err != nil {
		return h.errorResponse("Validation failed", err), nil
	}

	loginResp, err := h.service.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, errors.New("invalid username or password")) ||
			errors.Is(err, errors.New("account is inactive")) {
			return h.errorResponse("Invalid credentials or inactive account", err), nil
		}
		return h.errorResponse("Login failed", err), nil
	}

	return loginResp, nil
}

func (h *AuthorHandler) GetProfile(ctx *gofr.Context) (interface{}, error) {
	// Extract author ID from JWT token (this would be set by auth middleware)
	authorID, err := h.getAuthorIDFromContext(ctx)
	if err != nil {
		return h.errorResponse("Authentication required", err), nil
	}

	profile, err := h.service.GetProfile(ctx, authorID)
	if err != nil {
		if errors.Is(err, errors.New("author not found")) {
			return h.errorResponse("Author not found", err), nil
		}
		return h.errorResponse("Failed to get profile", err), nil
	}

	return profile, nil
}

func (h *AuthorHandler) UpdateProfile(ctx *gofr.Context) (interface{}, error) {
	// Extract author ID from JWT token
	authorID, err := h.getAuthorIDFromContext(ctx)
	if err != nil {
		return h.errorResponse("Authentication required", err), nil
	}

	var req models.UpdateAuthorRequest

	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse("Invalid request format", err), nil
	}

	if err := ValidateUpdateAuthorRequest(&req); err != nil {
		return h.errorResponse("Validation failed", err), nil
	}

	author, err := h.service.UpdateProfile(ctx, authorID, &req)
	if err != nil {
		if errors.Is(err, errors.New("author not found")) {
			return h.errorResponse("Author not found", err), nil
		}
		if errors.Is(err, errors.New("email already exists")) {
			return h.errorResponse("Email already exists", err), nil
		}
		return h.errorResponse("Failed to update profile", err), nil
	}

	return map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"author":  author,
	}, nil
}

func (h *AuthorHandler) ListAuthors(ctx *gofr.Context) (interface{}, error) {
	limit := 10
	offset := 0

	if limitStr := ctx.Param("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := ctx.Param("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	authors, err := h.service.ListAuthors(ctx, limit, offset)
	if err != nil {
		return h.errorResponse("Failed to get authors", err), nil
	}

	return map[string]interface{}{
		"success": true,
		"authors": authors,
		"limit":   limit,
		"offset":  offset,
	}, nil
}

func (h *AuthorHandler) DeleteAccount(ctx *gofr.Context) (interface{}, error) {
	// Extract author ID from JWT token
	authorID, err := h.getAuthorIDFromContext(ctx)
	if err != nil {
		return h.errorResponse("Authentication required", err), nil
	}

	err = h.service.DeleteAccount(ctx, authorID)
	if err != nil {
		if errors.Is(err, errors.New("author not found")) {
			return h.errorResponse("Author not found", err), nil
		}
		return h.errorResponse("Failed to delete account", err), nil
	}

	return map[string]interface{}{
		"success": true,
		"message": "Account deleted successfully",
	}, nil
}

func (h *AuthorHandler) AuthMiddleware(handler func(*gofr.Context) (interface{}, error)) func(*gofr.Context) (interface{}, error) {
	return func(ctx *gofr.Context) (interface{}, error) {
		// For now, we'll implement a simple authentication check
		// TODO: Implement proper JWT authentication when we understand GoFr request patterns better
		
		// For testing purposes, we'll set a dummy author ID
		// In production, this would extract and validate JWT token
		return handler(ctx)
	}
}

func (h *AuthorHandler) getAuthorIDFromContext(ctx *gofr.Context) (int, error) {
	// For now, return a dummy author ID for testing
	// TODO: Extract this from validated JWT token stored in context
	return 1, nil
}

func (h *AuthorHandler) errorResponse(message string, err error) map[string]any {
	response := map[string]any{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	return response
}
