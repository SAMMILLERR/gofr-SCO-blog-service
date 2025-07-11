package store

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"gofr-blog-service/models"

	"gofr.dev/pkg/gofr"
)

// Error definitions
var (
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrNotFound          = errors.New("record not found")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
	ErrInvalidID         = errors.New("invalid ID")
)

// PostStore handles database operations for posts
type PostStore struct{}

// NewPostStore creates a new post store instance
func NewPostStore() *PostStore {
	return &PostStore{}
}

// CreatePost persists a new blog post in the database
func (ps *PostStore) CreatePost(ctx *gofr.Context, post models.CreatePostRequest) (*models.Post, error) {
	var createdPost models.Post
	err := ctx.SQL.QueryRow(
		CreatePostQuery,
		post.Title, post.Content, post.Slug, post.AuthorID, post.Status,
	).Scan(
		&createdPost.ID, &createdPost.Title, &createdPost.Content, &createdPost.Slug,
		&createdPost.AuthorID, &createdPost.Status, &createdPost.CreatedAt, &createdPost.UpdatedAt,
	)

	if err != nil {
		return nil, errors.Join(ErrDatabaseOperation, err)
	}

	return &createdPost, nil
}

// GetPostByID retrieves a single post from the database by ID
func (ps *PostStore) GetPostByID(ctx *gofr.Context, id int) (*models.Post, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	var post models.Post
	err := ctx.SQL.QueryRow(GetPostByIDQuery, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID,
		&post.Status, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(ErrDatabaseOperation, err)
	}

	return &post, nil
}

// GetPosts retrieves posts from the database with pagination
func (ps *PostStore) GetPosts(ctx *gofr.Context, limit, offset int) ([]models.Post, error) {
	rows, err := ctx.SQL.Query(GetPostsQuery, limit, offset)
	if err != nil {
		return nil, errors.Join(ErrDatabaseOperation, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		scanErr := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID,
			&post.Status, &post.CreatedAt, &post.UpdatedAt,
		)
		if scanErr != nil {
			return nil, errors.Join(ErrDatabaseOperation, scanErr)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(ErrDatabaseOperation, err)
	}

	return posts, nil
}

// GetTotalPostCount returns the total number of posts in the database
func (ps *PostStore) GetTotalPostCount(ctx *gofr.Context) (int, error) {
	var totalCount int
	err := ctx.SQL.QueryRow(GetTotalPostCountQuery).Scan(&totalCount)
	if err != nil {
		return 0, errors.Join(ErrDatabaseOperation, err)
	}
	return totalCount, nil
}

// UpdatePost updates an existing post in the database
func (ps *PostStore) UpdatePost(ctx *gofr.Context, id int, req models.UpdatePostRequest) (*models.Post, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	// Build dynamic update query
	query, args := ps.buildUpdateQuery(id, req)
	if len(args) == 1 { // Only ID provided
		return nil, ErrNoFieldsToUpdate
	}

	var post models.Post
	err := ctx.SQL.QueryRow(query, args...).Scan(
		&post.ID, &post.Title, &post.Content, &post.Slug, &post.AuthorID,
		&post.Status, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(ErrDatabaseOperation, err)
	}

	return &post, nil
}

// DeletePost removes a post from the database by ID
func (ps *PostStore) DeletePost(ctx *gofr.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	result, err := ctx.SQL.Exec(DeletePostQuery, id)
	if err != nil {
		return errors.Join(ErrDatabaseOperation, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Join(ErrDatabaseOperation, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Helper method to build dynamic update queries
func (ps *PostStore) buildUpdateQuery(id int, req models.UpdatePostRequest) (query string, args []any) {
	setParts := []string{}
	args = []any{} // Using = instead of := since args is already declared in the return
	argIndex := 1

	if req.Title != "" {
		setParts = append(setParts, "title = $"+strconv.Itoa(argIndex))
		args = append(args, req.Title)
		argIndex++
	}
	if req.Content != "" {
		setParts = append(setParts, "content = $"+strconv.Itoa(argIndex))
		args = append(args, req.Content)
		argIndex++
	}
	if req.Slug != "" {
		setParts = append(setParts, "slug = $"+strconv.Itoa(argIndex))
		args = append(args, req.Slug)
		argIndex++
	}
	if req.Status != "" {
		setParts = append(setParts, "status = $"+strconv.Itoa(argIndex))
		args = append(args, req.Status)
		argIndex++
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	// Using = instead of := since query is already declared in the return
	query = "UPDATE posts SET " + strings.Join(setParts, ", ") +
		" WHERE id = $" + strconv.Itoa(argIndex) +
		" RETURNING id, title, content, slug, author_id, status, created_at, updated_at"

	return query, args
}
