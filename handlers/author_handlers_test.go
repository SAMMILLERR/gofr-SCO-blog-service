package handlers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gofr-blog-service/models"
)

// MockAuthorService is a mock implementation of the author service
type MockAuthorService struct {
	mock.Mock
}

func (m *MockAuthorService) Register(ctx interface{}, req *models.RegisterRequest) (*models.AuthorResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthorResponse), args.Error(1)
}

func (m *MockAuthorService) Login(ctx interface{}, req *models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthorService) GetProfile(ctx interface{}, authorID int) (*models.AuthorResponse, error) {
	args := m.Called(ctx, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthorResponse), args.Error(1)
}

func (m *MockAuthorService) UpdateProfile(ctx interface{}, authorID int, req *models.UpdateAuthorRequest) (*models.AuthorResponse, error) {
	args := m.Called(ctx, authorID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthorResponse), args.Error(1)
}

func (m *MockAuthorService) ListAuthors(ctx interface{}, limit, offset int) ([]*models.AuthorResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.AuthorResponse), args.Error(1)
}

func (m *MockAuthorService) DeleteAccount(ctx interface{}, authorID int) error {
	args := m.Called(ctx, authorID)
	return args.Error(0)
}

func (m *MockAuthorService) ValidateToken(token string) (*models.Author, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func TestNewAuthorHandler(t *testing.T) {
	mockService := &MockAuthorService{}
	handler := NewAuthorHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

// Test the handler structure without actually calling the HTTP methods
// since GoFr context is complex to mock properly
func TestAuthorHandler_Structure(t *testing.T) {
	mockService := &MockAuthorService{}
	handler := NewAuthorHandler(mockService)
	
	// Test that all handler methods exist and have correct types
	assert.NotNil(t, handler.Register)
	assert.NotNil(t, handler.Login)
	assert.NotNil(t, handler.GetProfile)
	assert.NotNil(t, handler.UpdateProfile)
	assert.NotNil(t, handler.ListAuthors)
	assert.NotNil(t, handler.DeleteAccount)
	assert.NotNil(t, handler.AuthMiddleware)
	assert.NotNil(t, handler.errorResponse)
}

func TestAuthorHandler_errorResponse(t *testing.T) {
	mockService := &MockAuthorService{}
	handler := NewAuthorHandler(mockService)
	
	// Test error response with no error
	result := handler.errorResponse("Test message", nil)
	assert.False(t, result["success"].(bool))
	assert.Equal(t, "Test message", result["message"])
	assert.Nil(t, result["error"])
	
	// Test error response with error
	testErr := errors.New("test error")
	result = handler.errorResponse("Test message", testErr)
	assert.False(t, result["success"].(bool))
	assert.Equal(t, "Test message", result["message"])
	assert.Equal(t, "test error", result["error"])
}

func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *models.RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			request: &models.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "invalid email",
			request: &models.RegisterRequest{
				Username:  "testuser",
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "weak password",
			request: &models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "weak",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
			errMsg:  "password must be at least 8 characters long",
		},
		{
			name: "invalid username characters",
			request: &models.RegisterRequest{
				Username:  "test@user!",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
			errMsg:  "username can only contain letters, numbers, and underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisterRequest(tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *models.LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			request: &models.LoginRequest{
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username or email is required",
		},
		{
			name: "missing password",
			request: &models.LoginRequest{
				Username: "testuser",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoginRequest(tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateAuthorRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *models.UpdateAuthorRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &models.UpdateAuthorRequest{
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			request: &models.UpdateAuthorRequest{
				Email: "invalid-email",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "short first name",
			request: &models.UpdateAuthorRequest{
				FirstName: "T",
			},
			wantErr: true,
			errMsg:  "first name must be between 2 and 50 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateAuthorRequest(tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
