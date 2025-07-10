package services

import (
	"fmt"

	"gofr-blog-service/models"
	"gofr-blog-service/store"

	"gofr.dev/pkg/gofr"
)

// PostService handles business logic for posts
type PostService struct {
	postStore *store.PostStore
}

// NewPostService creates a new post service instance
func NewPostService(postStore *store.PostStore) *PostService {
	return &PostService{
		postStore: postStore,
	}
}

// CreatePost creates a new blog post
func (ps *PostService) CreatePost(ctx *gofr.Context, req models.CreatePostRequest) (*models.Post, error) {
	// Let the handler handle validation
	post, err := ps.postStore.CreatePost(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	ctx.Logger.Infof("Post created successfully with ID: %d", post.ID)
	return post, nil
}

// GetPost retrieves a single post by ID
func (ps *PostService) GetPost(ctx *gofr.Context, id int) (*models.Post, error) {
	post, err := ps.postStore.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

// ListPosts retrieves posts with pagination
func (ps *PostService) ListPosts(ctx *gofr.Context, page, pageSize int) (*models.PostListResponse, error) {
	// Adjust pagination values if needed
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize

	// Get total count from store
	totalCount, err := ps.postStore.GetTotalPostCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts from store
	posts, err := ps.postStore.GetPosts(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	return &models.PostListResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdatePost updates an existing post
func (ps *PostService) UpdatePost(ctx *gofr.Context, id int, req models.UpdatePostRequest) (*models.Post, error) {
	// Let the handler handle validation of id
	post, err := ps.postStore.UpdatePost(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	ctx.Logger.Infof("Post updated successfully: %d", post.ID)
	return post, nil
}

// DeletePost removes a post by ID
func (ps *PostService) DeletePost(ctx *gofr.Context, id int) error {
	// Let the handler handle validation of id
	err := ps.postStore.DeletePost(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	ctx.Logger.Infof("Post deleted successfully: %d", id)
	return nil
}
