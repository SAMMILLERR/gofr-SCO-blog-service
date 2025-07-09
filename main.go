package main

import (
	"gofr-blog-service/handlers"
	"gofr-blog-service/services"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/migration"
)

func main() {
	app := gofr.New()

	// Add database migrations
	app.Migrate(map[int64]migration.Migrate{
		1: createPostsTable(),
	})

	// Initialize services
	postService := services.NewPostService()

	// Initialize handlers
	postHandler := handlers.NewPostHandler(postService)

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (any, error) {
		return map[string]string{
			"status":  "healthy",
			"service": "gofr-blog-service",
			"version": "1.0.0",
		}, nil
	})

	// Simplified Post routes (removed /api/v1 prefix)
	app.GET("/posts", postHandler.ListPosts)
	app.GET("/posts/{id}", postHandler.GetPost)
	app.POST("/posts", postHandler.CreatePost)
	app.PUT("/posts/{id}", postHandler.UpdatePost)
	app.DELETE("/posts/{id}", postHandler.DeletePost)

	app.Run()
}

// createPostsTable defines the migration for creating posts table
func createPostsTable() migration.Migrate {
	return migration.Migrate{
		UP: func(datasource migration.Datasource) error {
			_, err := datasource.SQL.Exec(`
				CREATE TABLE IF NOT EXISTS posts (
					id SERIAL PRIMARY KEY,
					title VARCHAR(200) NOT NULL,
					content TEXT NOT NULL,
					slug VARCHAR(200) NOT NULL UNIQUE,
					author_id INTEGER NOT NULL,
					status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
					created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
				);

				-- Create indexes for better performance
				CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
				CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
				CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
				CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
			`)
			return err
		},
	}
}
