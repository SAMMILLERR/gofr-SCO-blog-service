package main

import (
	// Standard library imports first (none here)

	// External packages next
	"gofr.dev/pkg/gofr"

	// Internal packages last, sorted alphabetically
	"gofr-blog-service/handlers"
	"gofr-blog-service/migrations"
	"gofr-blog-service/services"
	"gofr-blog-service/store"
)

func main() {
	app := gofr.New()

	// Add database migrations from migrations package
	app.Migrate(migrations.GetMigrations())

	// Initialize store (new layer)
	postStore := store.NewPostStore()

	// Initialize services with store dependency
	postService := services.NewPostService(postStore)

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

	// Simplified Post routes
	app.GET("/posts", postHandler.ListPosts)
	app.GET("/posts/{id}", postHandler.GetPost)
	app.POST("/posts", postHandler.CreatePost)
	app.PUT("/posts/{id}", postHandler.UpdatePost)
	app.DELETE("/posts/{id}", postHandler.DeletePost)

	app.Run()
}
