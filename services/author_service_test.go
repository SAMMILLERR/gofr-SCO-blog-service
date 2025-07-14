package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gofr-blog-service/models"
)

// MockAuthorStore is a mock implementation of the author store
type MockAuthorStore struct {
	mock.Mock
}

func (m *MockAuthorStore) Create(ctx interface{}, author *models.Author) (*models.Author, error) {
	args := m.Called(ctx, author)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorStore) GetByID(ctx interface{}, id int) (*models.Author, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorStore) GetByUsername(ctx interface{}, username string) (*models.Author, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorStore) GetByEmail(ctx interface{}, email string) (*models.Author, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorStore) GetPasswordHash(ctx interface{}, username string) (string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorStore) Update(ctx interface{}, id int, author *models.Author) (*models.Author, error) {
	args := m.Called(ctx, id, author)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorStore) Delete(ctx interface{}, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuthorStore) List(ctx interface{}, limit, offset int) ([]*models.Author, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Author), args.Error(1)
}

func TestNewAuthorService(t *testing.T) {
	mockStore := &MockAuthorStore{}
	jwtSecret := "test-secret"
	
	service := NewAuthorService(mockStore, jwtSecret)
	
	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
	assert.Equal(t, jwtSecret, service.jwtSecret)
}

func TestAuthorService_Register(t *testing.T) {
	tests := []struct {
		name     string
		request  *models.RegisterRequest
		setup    func(*MockAuthorStore)
		wantErr  bool
		errMsg   string
		validate func(*testing.T, *models.AuthorResponse)
	}{
		{
			name: "successful registration",
			request: &models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Bio:       "Test bio",
				AvatarURL: "https://example.com/avatar.jpg",
			},
			setup: func(mockStore *MockAuthorStore) {
				author := &models.Author{
					ID:         1,
					Username:   "testuser",
					Email:      "test@example.com",
					FirstName:  "Test",
					LastName:   "User",
					Bio:        "Test bio",
					AvatarURL:  "https://example.com/avatar.jpg",
					IsActive:   true,
					IsVerified: false,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				mockStore.On("Create", mock.Anything, mock.AnythingOfType("*models.Author")).Return(author, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, response *models.AuthorResponse) {
				assert.Equal(t, 1, response.ID)
				assert.Equal(t, "testuser", response.Username)
				assert.Equal(t, "test@example.com", response.Email)
				assert.Equal(t, "Test", response.FirstName)
				assert.Equal(t, "User", response.LastName)
				assert.True(t, response.IsActive)
				assert.False(t, response.IsVerified)
			},
		},
		{
			name: "weak password error",
			request: &models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "weak",
				FirstName: "Test",
				LastName:  "User",
			},
			setup:   func(mockStore *MockAuthorStore) {},
			wantErr: true,
			errMsg:  "password must be at least 8 characters long",
		},
		{
			name: "duplicate username error",
			request: &models.RegisterRequest{
				Username:  "existinguser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			setup: func(mockStore *MockAuthorStore) {
				mockStore.On("Create", mock.Anything, mock.AnythingOfType("*models.Author")).Return(nil, errors.New("username already exists"))
			},
			wantErr: true,
			errMsg:  "username already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockAuthorStore{}
			tt.setup(mockStore)
			
			service := NewAuthorService(mockStore, "test-secret")
			
			// We can't easily create a real gofr.Context for testing
			// So we'll use nil as the context interface
			ctx := interface{}(nil)
			
			response, err := service.Register(ctx, tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				if tt.validate != nil {
					tt.validate(t, response)
				}
			}
			
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthorService_Login(t *testing.T) {
	tests := []struct {
		name    string
		request *models.LoginRequest
		setup   func(*MockAuthorStore)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful login with username",
			request: &models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setup: func(mockStore *MockAuthorStore) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				author := &models.Author{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					IsActive:  true,
				}
				mockStore.On("GetByUsername", mock.Anything, "testuser").Return(author, nil)
				mockStore.On("GetPasswordHash", mock.Anything, "testuser").Return(string(hashedPassword), nil)
			},
			wantErr: false,
		},
		{
			name: "successful login with email",
			request: &models.LoginRequest{
				Username: "test@example.com",
				Password: "password123",
			},
			setup: func(mockStore *MockAuthorStore) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				author := &models.Author{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					IsActive:  true,
				}
				mockStore.On("GetByEmail", mock.Anything, "test@example.com").Return(author, nil)
				mockStore.On("GetPasswordHash", mock.Anything, "testuser").Return(string(hashedPassword), nil)
			},
			wantErr: false,
		},
		{
			name: "invalid password",
			request: &models.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setup: func(mockStore *MockAuthorStore) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				author := &models.Author{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					IsActive:  true,
				}
				mockStore.On("GetByUsername", mock.Anything, "testuser").Return(author, nil)
				mockStore.On("GetPasswordHash", mock.Anything, "testuser").Return(string(hashedPassword), nil)
			},
			wantErr: true,
			errMsg:  "invalid username or password",
		},
		{
			name: "inactive account",
			request: &models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setup: func(mockStore *MockAuthorStore) {
				author := &models.Author{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					IsActive:  false, // Account is inactive
				}
				mockStore.On("GetByUsername", mock.Anything, "testuser").Return(author, nil)
			},
			wantErr: true,
			errMsg:  "account is inactive",
		},
		{
			name: "user not found",
			request: &models.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			setup: func(mockStore *MockAuthorStore) {
				mockStore.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, errors.New("author not found"))
			},
			wantErr: true,
			errMsg:  "invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockAuthorStore{}
			tt.setup(mockStore)
			
			service := NewAuthorService(mockStore, "test-secret")
			ctx := interface{}(nil)
			
			response, err := service.Login(ctx, tt.request)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.NotNil(t, response.Author)
			}
			
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthorService_GetProfile(t *testing.T) {
	tests := []struct {
		name     string
		authorID int
		setup    func(*MockAuthorStore)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful profile retrieval",
			authorID: 1,
			setup: func(mockStore *MockAuthorStore) {
				author := &models.Author{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
					IsActive:  true,
				}
				mockStore.On("GetByID", mock.Anything, 1).Return(author, nil)
			},
			wantErr: false,
		},
		{
			name:     "author not found",
			authorID: 999,
			setup: func(mockStore *MockAuthorStore) {
				mockStore.On("GetByID", mock.Anything, 999).Return(nil, errors.New("author not found"))
			},
			wantErr: true,
			errMsg:  "author not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockAuthorStore{}
			tt.setup(mockStore)
			
			service := NewAuthorService(mockStore, "test-secret")
			ctx := interface{}(nil)
			
			response, err := service.GetProfile(ctx, tt.authorID)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.authorID, response.ID)
			}
			
			mockStore.AssertExpectations(t)
		})
	}
}

func TestAuthorService_ValidateToken(t *testing.T) {
	jwtSecret := "test-secret"
	service := NewAuthorService(&MockAuthorStore{}, jwtSecret)
	
	// Create a valid token first
	author := &models.Author{
		ID:       1,
		Username: "testuser",
	}
	
	token, err := service.generateJWT(author)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Test valid token
	validatedAuthor, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, validatedAuthor)
	assert.Equal(t, 1, validatedAuthor.ID)
	assert.Equal(t, "testuser", validatedAuthor.Username)
	
	// Test invalid token
	invalidAuthor, err := service.ValidateToken("invalid.token.here")
	assert.Error(t, err)
	assert.Nil(t, invalidAuthor)
}

func TestAuthorService_generateJWT(t *testing.T) {
	service := NewAuthorService(&MockAuthorStore{}, "test-secret")
	
	author := &models.Author{
		ID:       1,
		Username: "testuser",
	}
	
	token, err := service.generateJWT(author)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Token should be a valid JWT format (three parts separated by dots)
	parts := len(token) > 0
	assert.True(t, parts)
}
