package main

import (
	"gofr.dev/pkg/gofr"

	"gofr-blog-service/handlers"
	"gofr-blog-service/migrations"
	"gofr-blog-service/services"
	"gofr-blog-service/store"
)

func main() {
	app := gofr.New()

	// Add database migrations from migrations package
	app.Migrate(migrations.All())

	// JWT secret (in production, this should come from environment variables)
	jwtSecret := app.Config.Get("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// Initialize stores
	postStore := store.NewPostStore()
	authorStore := store.NewAuthorStore()

	// Initialize services with store dependencies
	postService := services.NewPostService(postStore)
	authorService := services.NewAuthorService(authorStore, jwtSecret)

	// Initialize handlers
	postHandler := handlers.NewPostHandler(postService)
	authorHandler := handlers.NewAuthorHandler(authorService)

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (any, error) {
		return map[string]string{
			"status":  "healthy",
			"service": "gofr-blog-service",
			"version": "1.0.0",
		}, nil
	})

	// Public author routes
	app.POST("/auth/register", authorHandler.Register)
	app.POST("/auth/login", authorHandler.Login)
	app.GET("/authors", authorHandler.ListAuthors)

	// Protected author routes (require authentication)
	app.GET("/auth/profile", authorHandler.AuthMiddleware(authorHandler.GetProfile))
	app.PUT("/auth/profile", authorHandler.AuthMiddleware(authorHandler.UpdateProfile))
	app.DELETE("/auth/account", authorHandler.AuthMiddleware(authorHandler.DeleteAccount))

	// Post routes (some may require authentication in the future)
	app.GET("/posts", postHandler.ListPosts)
	app.GET("/posts/{id}", postHandler.GetPost)
	app.POST("/posts", authorHandler.AuthMiddleware(postHandler.CreatePost))
	app.PUT("/posts/{id}", authorHandler.AuthMiddleware(postHandler.UpdatePost))
	app.DELETE("/posts/{id}", authorHandler.AuthMiddleware(postHandler.DeletePost))

	app.Run()
}
