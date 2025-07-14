package store

// SQL queries for post store operations
const (
	// CreatePostQuery inserts a new post into the database
	CreatePostQuery = `
		INSERT INTO posts (title, content, slug, author_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, title, content, slug, author_id, status, created_at, updated_at
	`

	// GetPostByIDQuery retrieves a post by its ID
	GetPostByIDQuery = `
		SELECT id, title, content, slug, author_id, status, created_at, updated_at
		FROM posts WHERE id = $1
	`

	// GetPostsQuery retrieves posts with pagination
	GetPostsQuery = `
		SELECT id, title, content, slug, author_id, status, created_at, updated_at
		FROM posts 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	// GetTotalPostCountQuery counts the total number of posts
	GetTotalPostCountQuery = `SELECT COUNT(*) FROM posts`

	// DeletePostQuery deletes a post by ID
	DeletePostQuery = `DELETE FROM posts WHERE id = $1`

	// UpdatePostBaseQuery is the base for dynamic update queries
	UpdatePostBaseQuery = `UPDATE posts SET %s, updated_at = NOW() WHERE id = $%d 
		RETURNING id, title, content, slug, author_id, status, created_at, updated_at`
)
