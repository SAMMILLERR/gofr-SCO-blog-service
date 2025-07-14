package store

import (
	"gofr-blog-service/models"

	"github.com/stretchr/testify/mock"
	"gofr.dev/pkg/gofr"
)

// MockPostStore is a mock implementation of the post store for testing
type MockPostStore struct {
	mock.Mock
}

// CreatePost mocks the CreatePost method
func (m *MockPostStore) CreatePost(ctx *gofr.Context, post models.CreatePostRequest) (*models.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

// GetPostByID mocks the GetPostByID method
func (m *MockPostStore) GetPostByID(ctx *gofr.Context, id int) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

// GetPosts mocks the GetPosts method
func (m *MockPostStore) GetPosts(ctx *gofr.Context, limit, offset int) ([]models.Post, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Post), args.Error(1)
}

// GetTotalPostCount mocks the GetTotalPostCount method
func (m *MockPostStore) GetTotalPostCount(ctx *gofr.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// UpdatePost mocks the UpdatePost method
func (m *MockPostStore) UpdatePost(ctx *gofr.Context, id int, req models.UpdatePostRequest) (*models.Post, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

// DeletePost mocks the DeletePost method
func (m *MockPostStore) DeletePost(ctx *gofr.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
