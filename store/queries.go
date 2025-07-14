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

	// Author queries
	createAuthorQuery = `
		INSERT INTO authors (username, email, password_hash, first_name, last_name, bio, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, username, email, first_name, last_name, bio, avatar_url, is_active, is_verified, created_at, updated_at
	`

	getAuthorByIDQuery = `
		SELECT id, username, email, first_name, last_name, bio, avatar_url, is_active, is_verified, created_at, updated_at
		FROM authors WHERE id = $1
	`

	getAuthorByUsernameQuery = `
		SELECT id, username, email, first_name, last_name, bio, avatar_url, is_active, is_verified, created_at, updated_at
		FROM authors WHERE username = $1
	`

	getAuthorByEmailQuery = `
		SELECT id, username, email, first_name, last_name, bio, avatar_url, is_active, is_verified, created_at, updated_at
		FROM authors WHERE email = $1
	`

	getPasswordHashQuery = `
		SELECT password_hash FROM authors WHERE username = $1 AND is_active = true
	`

	updateAuthorQuery = `
		UPDATE authors 
		SET email = $1, first_name = $2, last_name = $3, bio = $4, avatar_url = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, username, email, first_name, last_name, bio, avatar_url, is_active, is_verified, created_at, updated_at
	`

	deleteAuthorQuery = `
		DELETE FROM authors WHERE id = $1
	`

	listAuthorsQuery = `
		SELECT id, username, email, first_name, last_name, bio, avatar_url, is_active, is_verified, created_at, updated_at
		FROM authors 
		WHERE is_active = true
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
)
