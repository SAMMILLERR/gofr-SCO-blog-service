package main

import (
	"gofr-blog-service/internal/handlers"
	"gofr-blog-service/internal/services"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

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

	// Post routes
	app.GET("/api/v1/posts", postHandler.ListPosts)
	app.GET("/api/v1/posts/{id}", postHandler.GetPost)
	app.POST("/api/v1/posts", postHandler.CreatePost)
	app.PUT("/api/v1/posts/{id}", postHandler.UpdatePost)
	app.DELETE("/api/v1/posts/{id}", postHandler.DeletePost)

	app.Run()
}
