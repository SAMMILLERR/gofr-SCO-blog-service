package services

import (
	"testing"

	"gofr-blog-service/store"
)

// TestNewPostService tests the creation of a new post service
func TestNewPostService(t *testing.T) {
	// Create a mock store
	mockStore := &store.PostStore{}

	// Create a new service
	service := NewPostService(mockStore)

	// Check that the service has the correct store
	if service.postStore != mockStore {
		t.Errorf("Expected postStore to be %v, got %v", mockStore, service.postStore)
	}
}
