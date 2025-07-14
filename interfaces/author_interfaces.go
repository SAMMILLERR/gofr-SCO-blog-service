package interfaces

import (
	"gofr-blog-service/models"
)

// AuthorStoreInterface defines the contract for author data operations
type AuthorStoreInterface interface {
	Create(ctx interface{}, author *models.Author) (*models.Author, error)
	GetByID(ctx interface{}, id int) (*models.Author, error)
	GetByUsername(ctx interface{}, username string) (*models.Author, error)
	GetByEmail(ctx interface{}, email string) (*models.Author, error)
	GetPasswordHash(ctx interface{}, username string) (string, error)
	Update(ctx interface{}, id int, author *models.Author) (*models.Author, error)
	Delete(ctx interface{}, id int) error
	List(ctx interface{}, limit, offset int) ([]*models.Author, error)
}

// AuthorServiceInterface defines the contract for author business logic
type AuthorServiceInterface interface {
	Register(ctx interface{}, req *models.RegisterRequest) (*models.AuthorResponse, error)
	Login(ctx interface{}, req *models.LoginRequest) (*models.LoginResponse, error)
	GetProfile(ctx interface{}, authorID int) (*models.AuthorResponse, error)
	UpdateProfile(ctx interface{}, authorID int, req *models.UpdateAuthorRequest) (*models.AuthorResponse, error)
	ListAuthors(ctx interface{}, limit, offset int) ([]*models.AuthorResponse, error)
	DeleteAccount(ctx interface{}, authorID int) error
	ValidateToken(token string) (*models.Author, error)
}
