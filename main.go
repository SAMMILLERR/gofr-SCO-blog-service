package main

import (
	"gofr-blog-service/database"
	"gofr-blog-service/handlers"
	"gofr-blog-service/services"
	"gofr-blog-service/store"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// Add database migrations - moved to database package
	app.Migrate(database.GetMigrations())

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
