package services

import (
    "fmt"
    "strings"
    "gofr-blog-service/internal/models"
    "gofr.dev/pkg/gofr"
)

// PostService handles business logic for posts
type PostService struct {
    // Using dependency injection pattern - no direct DB dependency
}

// NewPostService creates a new post service instance
func NewPostService() *PostService {
    return &PostService{}
}

// CreatePost creates a new blog post with validation
func (ps *PostService) CreatePost(ctx *gofr.Context, req models.CreatePostRequest) (*models.Post, error) {
    // Input validation decorator pattern
    if err := ps.validateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Database operation with error handling decorator
    post, err := ps.createPostInDB(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create post: %w", err)
    }

    // Success logging decorator could be added here
    ctx.Logger.Infof("Post created successfully with ID: %d", post.ID)
    
    return post, nil
}

// GetPost retrieves a single post by ID
func (ps *PostService) GetPost(ctx *gofr.Context, id int) (*models.Post, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid post ID: %d", id)
    }

    query := `
        SELECT id, title, content, slug, author_id, status, created_at, updated_at
        FROM posts WHERE id = $1
    `
    
    var post models.Post
    err := ctx.SQL.QueryRow(query, id).
        Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("post not found: %w", err)
    }

    return &post, nil
}

// ListPosts retrieves posts with pagination
func (ps *PostService) ListPosts(ctx *gofr.Context, page, pageSize int) (*models.PostListResponse, error) {
    // Pagination validation decorator
    page, pageSize = ps.validatePagination(page, pageSize)
    offset := (page - 1) * pageSize
    
    // Get total count
    totalCount, err := ps.getTotalPostCount(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to count posts: %w", err)
    }

    // Get posts with error handling
    posts, err := ps.getPostsFromDB(ctx, pageSize, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to query posts: %w", err)
    }

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
    if id <= 0 {
        return nil, fmt.Errorf("invalid post ID: %d", id)
    }

    // Build dynamic update query (decorator pattern for query building)
    query, args := ps.buildUpdateQuery(id, req)
    if len(args) == 1 { // Only ID provided
        return nil, fmt.Errorf("no fields to update")
    }

    var post models.Post
    err := ctx.SQL.QueryRow(query, args...).
        Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("failed to update post: %w", err)
    }

    ctx.Logger.Infof("Post updated successfully: %d", post.ID)
    return &post, nil
}

// DeletePost removes a post by ID
func (ps *PostService) DeletePost(ctx *gofr.Context, id int) error {
    if id <= 0 {
        return fmt.Errorf("invalid post ID: %d", id)
    }

    query := "DELETE FROM posts WHERE id = $1"
    result, err := ctx.SQL.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete post: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("post not found")
    }

    ctx.Logger.Infof("Post deleted successfully: %d", id)
    return nil
}

// Private helper methods (decorator pattern for separation of concerns)

func (ps *PostService) validateCreateRequest(req models.CreatePostRequest) error {
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

func (ps *PostService) validatePagination(page, pageSize int) (int, int) {
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 || pageSize > 100 {
        pageSize = 10
    }
    return page, pageSize
}

func (ps *PostService) createPostInDB(ctx *gofr.Context, req models.CreatePostRequest) (*models.Post, error) {
    query := `
        INSERT INTO posts (title, content, slug, author_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, title, content, slug, author_id, status, created_at, updated_at
    `
    
    var post models.Post
    err := ctx.SQL.QueryRow(query, req.Title, req.Content, req.Slug, req.AuthorID, req.Status).
        Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt)
    
    return &post, err
}

func (ps *PostService) getTotalPostCount(ctx *gofr.Context) (int, error) {
    var totalCount int
    countQuery := "SELECT COUNT(*) FROM posts"
    err := ctx.SQL.QueryRow(countQuery).Scan(&totalCount)
    return totalCount, err
}

func (ps *PostService) getPostsFromDB(ctx *gofr.Context, limit, offset int) ([]models.Post, error) {
    query := `
        SELECT id, title, content, slug, author_id, status, created_at, updated_at
        FROM posts 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    `
    
    rows, err := ctx.SQL.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

func (ps *PostService) buildUpdateQuery(id int, req models.UpdatePostRequest) (string, []any) {
    setParts := []string{}
    args := []any{}
    argIndex := 1

    if req.Title != "" {
        setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
        args = append(args, req.Title)
        argIndex++
    }
    if req.Content != "" {
        setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
        args = append(args, req.Content)
        argIndex++
    }
    if req.Slug != "" {
        setParts = append(setParts, fmt.Sprintf("slug = $%d", argIndex))
        args = append(args, req.Slug)
        argIndex++
    }
    if req.Status != "" {
        setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
        args = append(args, req.Status)
        argIndex++
    }

    setParts = append(setParts, "updated_at = NOW()")
    args = append(args, id)

    query := fmt.Sprintf(`
        UPDATE posts 
        SET %s 
        WHERE id = $%d
        RETURNING id, title, content, slug, author_id, status, created_at, updated_at
    `, strings.Join(setParts, ", "), argIndex)

    return query, args
}
