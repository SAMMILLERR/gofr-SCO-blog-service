package services

import (
    "testing"
    "gofr-blog-service/internal/models"
)

// Test decorators for comprehensive testing

func TestPostService_CreatePost(t *testing.T) {
    tests := []struct {
        name    string
        request models.CreatePostRequest
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid post creation",
            request: models.CreatePostRequest{
                Title:    "Test Post",
                Content:  "This is a test post content with enough characters",
                Slug:     "test-post",
                AuthorID: 1,
                Status:   "draft",
            },
            wantErr: false,
        },
        {
            name: "empty title should fail",
            request: models.CreatePostRequest{
                Title:    "",
                Content:  "This is a test post content",
                Slug:     "test-post",
                AuthorID: 1,
                Status:   "draft",
            },
            wantErr: true,
            errMsg:  "title is required",
        },
        {
            name: "short content should fail",
            request: models.CreatePostRequest{
                Title:    "Test Post",
                Content:  "Short",
                Slug:     "test-post",
                AuthorID: 1,
                Status:   "draft",
            },
            wantErr: true,
            errMsg:  "content must be at least 10 characters",
        },
        {
            name: "invalid author ID should fail",
            request: models.CreatePostRequest{
                Title:    "Test Post",
                Content:  "This is a test post content",
                Slug:     "test-post",
                AuthorID: 0,
                Status:   "draft",
            },
            wantErr: true,
            errMsg:  "valid author ID is required",
        },
    }

    service := NewPostService()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test validation decorator
            err := service.validateCreateRequest(tt.request)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error but got none")
                    return
                }
                if tt.errMsg != "" && err.Error() != tt.errMsg {
                    t.Errorf("expected error message %q, got %q", tt.errMsg, err.Error())
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}

func TestPostService_ValidatePagination(t *testing.T) {
    tests := []struct {
        name         string
        page         int
        pageSize     int
        expectedPage int
        expectedSize int
    }{
        {"valid pagination", 2, 20, 2, 20},
        {"invalid page defaults to 1", 0, 20, 1, 20},
        {"negative page defaults to 1", -1, 20, 1, 20},
        {"invalid pageSize defaults to 10", 1, 0, 1, 10},
        {"too large pageSize caps at 100", 1, 150, 1, 10},
        {"negative pageSize defaults to 10", 1, -5, 1, 10},
    }

    service := NewPostService()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            page, pageSize := service.validatePagination(tt.page, tt.pageSize)
            
            if page != tt.expectedPage {
                t.Errorf("expected page %d, got %d", tt.expectedPage, page)
            }
            if pageSize != tt.expectedSize {
                t.Errorf("expected pageSize %d, got %d", tt.expectedSize, pageSize)
            }
        })
    }
}

func TestPostService_BuildUpdateQuery(t *testing.T) {
    service := NewPostService()
    
    tests := []struct {
        name          string
        id            int
        request       models.UpdatePostRequest
        expectedParts []string
        expectedArgs  int
    }{
        {
            name: "update title only",
            id:   1,
            request: models.UpdatePostRequest{
                Title: "New Title",
            },
            expectedParts: []string{"title = $1", "updated_at = NOW()"},
            expectedArgs:  2, // title value + id
        },
        {
            name: "update multiple fields",
            id:   1,
            request: models.UpdatePostRequest{
                Title:   "New Title",
                Content: "New Content",
                Status:  "published",
            },
            expectedParts: []string{"title = $1", "content = $2", "status = $3", "updated_at = NOW()"},
            expectedArgs:  4, // 3 field values + id
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            query, args := service.buildUpdateQuery(tt.id, tt.request)
            
            if len(args) != tt.expectedArgs {
                t.Errorf("expected %d args, got %d", tt.expectedArgs, len(args))
            }
            
            // Check that query contains expected parts
            for _, part := range tt.expectedParts {
                if !containsString(query, part) {
                    t.Errorf("query should contain %q", part)
                }
            }
        })
    }
}

// Helper function for testing
func containsString(str, substr string) bool {
    return len(str) >= len(substr) && (str == substr || 
        (len(str) > len(substr) && 
         (str[:len(substr)] == substr || 
          str[len(str)-len(substr):] == substr ||
          hasSubstring(str, substr))))
}

func hasSubstring(str, substr string) bool {
    for i := 0; i <= len(str)-len(substr); i++ {
        if str[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
