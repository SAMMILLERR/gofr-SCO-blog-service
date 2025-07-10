package store

import (
	"fmt"
	"strings"

	"gofr-blog-service/models"

	"gofr.dev/pkg/gofr"
)

// PostStore handles database operations for posts
type PostStore struct{}

// NewPostStore creates a new post store instance
func NewPostStore() *PostStore {
	return &PostStore{}
}

// CreatePost persists a new blog post in the database
func (ps *PostStore) CreatePost(ctx *gofr.Context, post models.CreatePostRequest) (*models.Post, error) {
	query := `
        INSERT INTO posts (title, content, slug, author_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, title, content, slug, author_id, status, created_at, updated_at
    `

	var createdPost models.Post
	err := ctx.SQL.QueryRow(query, post.Title, post.Content, post.Slug, post.AuthorID, post.Status).
		Scan(&createdPost.ID, &createdPost.Title, &createdPost.Content, &createdPost.Slug, &createdPost.AuthorID,
			&createdPost.Status, &createdPost.CreatedAt, &createdPost.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &createdPost, nil
}

// GetPostByID retrieves a single post from the database by ID
func (ps *PostStore) GetPostByID(ctx *gofr.Context, id int) (*models.Post, error) {
	query := `
        SELECT id, title, content, slug, author_id, status, created_at, updated_at
        FROM posts WHERE id = $1
    `

	var post models.Post
	err := ctx.SQL.QueryRow(query, id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID,
			&post.Status, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &post, nil
}

// GetPosts retrieves posts from the database with pagination
func (ps *PostStore) GetPosts(ctx *gofr.Context, limit, offset int) ([]models.Post, error) {
	query := `
        SELECT id, title, content, slug, author_id, status, created_at, updated_at
        FROM posts 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    `

	rows, err := ctx.SQL.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID,
			&post.Status, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return posts, nil
}

// GetTotalPostCount returns the total number of posts in the database
func (ps *PostStore) GetTotalPostCount(ctx *gofr.Context) (int, error) {
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM posts"
	err := ctx.SQL.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("count query error: %w", err)
	}
	return totalCount, nil
}

// UpdatePost updates an existing post in the database
func (ps *PostStore) UpdatePost(ctx *gofr.Context, id int, req models.UpdatePostRequest) (*models.Post, error) {
	// Build dynamic update query
	query, args := ps.buildUpdateQuery(id, req)
	if len(args) == 1 { // Only ID provided
		return nil, fmt.Errorf("no fields to update")
	}

	var post models.Post
	err := ctx.SQL.QueryRow(query, args...).
		Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID,
			&post.Status, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("update query error: %w", err)
	}

	return &post, nil
}

// DeletePost removes a post from the database by ID
func (ps *PostStore) DeletePost(ctx *gofr.Context, id int) error {
	query := "DELETE FROM posts WHERE id = $1"
	result, err := ctx.SQL.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete query error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// Helper method to build dynamic update queries
func (ps *PostStore) buildUpdateQuery(id int, req models.UpdatePostRequest) (query string, args []any) {
	setParts := []string{}
	args = []any{}
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

	query = fmt.Sprintf(`
        UPDATE posts 
        SET %s 
        WHERE id = $%d
        RETURNING id, title, content, slug, author_id, status, created_at, updated_at
    `, strings.Join(setParts, ", "), argIndex)

	return query, args
}
