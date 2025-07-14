package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"

	"gofr-blog-service/interfaces"
	"gofr-blog-service/models"
)

var (
	errInvalidCredentials = errors.New("invalid username or password")
	errWeakPassword      = errors.New("password must be at least 8 characters long")
	errAccountInactive   = errors.New("account is inactive")
)

type AuthorService struct {
	store     interfaces.AuthorStoreInterface
	jwtSecret string
}

func NewAuthorService(store interfaces.AuthorStoreInterface, jwtSecret string) *AuthorService {
	return &AuthorService{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthorService) Register(ctx interface{}, req *models.RegisterRequest) (*models.AuthorResponse, error) {
	// Validate password strength
	if len(req.Password) < 8 {
		return nil, errWeakPassword
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create author model
	author := &models.Author{
		Username:     strings.ToLower(strings.TrimSpace(req.Username)),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(hashedPassword),
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Bio:          req.Bio,
		AvatarURL:    req.AvatarURL,
		IsActive:     true,
		IsVerified:   false,
	}

	// Create the author
	createdAuthor, err := s.store.Create(ctx, author)
	if err != nil {
		return nil, err
	}

	// Convert to response
	return &models.AuthorResponse{
		ID:         createdAuthor.ID,
		Username:   createdAuthor.Username,
		Email:      createdAuthor.Email,
		FirstName:  createdAuthor.FirstName,
		LastName:   createdAuthor.LastName,
		Bio:        createdAuthor.Bio,
		AvatarURL:  createdAuthor.AvatarURL,
		IsActive:   createdAuthor.IsActive,
		IsVerified: createdAuthor.IsVerified,
		CreatedAt:  createdAuthor.CreatedAt,
		UpdatedAt:  createdAuthor.UpdatedAt,
	}, nil
}

func (s *AuthorService) Login(ctx interface{}, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get author by username or email
	var author *models.Author
	var err error

	if strings.Contains(req.Username, "@") {
		author, err = s.store.GetByEmail(ctx, strings.ToLower(req.Username))
	} else {
		author, err = s.store.GetByUsername(ctx, strings.ToLower(req.Username))
	}

	if err != nil {
		return nil, errInvalidCredentials
	}

	// Check if account is active
	if !author.IsActive {
		return nil, errAccountInactive
	}

	// Get password hash
	passwordHash, err := s.store.GetPasswordHash(ctx, author.Username)
	if err != nil {
		return nil, errInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, errInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateJWT(author)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		Author: models.AuthorResponse{
			ID:         author.ID,
			Username:   author.Username,
			Email:      author.Email,
			FirstName:  author.FirstName,
			LastName:   author.LastName,
			Bio:        author.Bio,
			AvatarURL:  author.AvatarURL,
			IsActive:   author.IsActive,
			IsVerified: author.IsVerified,
			CreatedAt:  author.CreatedAt,
			UpdatedAt:  author.UpdatedAt,
		},
	}, nil
}

func (s *AuthorService) GetProfile(ctx interface{}, authorID int) (*models.AuthorResponse, error) {
	author, err := s.store.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return &models.AuthorResponse{
		ID:         author.ID,
		Username:   author.Username,
		Email:      author.Email,
		FirstName:  author.FirstName,
		LastName:   author.LastName,
		Bio:        author.Bio,
		AvatarURL:  author.AvatarURL,
		IsActive:   author.IsActive,
		IsVerified: author.IsVerified,
		CreatedAt:  author.CreatedAt,
		UpdatedAt:  author.UpdatedAt,
	}, nil
}

func (s *AuthorService) UpdateProfile(ctx interface{}, authorID int, req *models.UpdateAuthorRequest) (*models.AuthorResponse, error) {
	// Trim and validate input
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	// Convert UpdateAuthorRequest to Author model
	author := &models.Author{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	
	// Handle pointer fields
	if req.Bio != nil {
		author.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		author.AvatarURL = *req.AvatarURL
	}

	updatedAuthor, err := s.store.Update(ctx, authorID, author)
	if err != nil {
		return nil, err
	}

	return &models.AuthorResponse{
		ID:         updatedAuthor.ID,
		Username:   updatedAuthor.Username,
		Email:      updatedAuthor.Email,
		FirstName:  updatedAuthor.FirstName,
		LastName:   updatedAuthor.LastName,
		Bio:        updatedAuthor.Bio,
		AvatarURL:  updatedAuthor.AvatarURL,
		IsActive:   updatedAuthor.IsActive,
		IsVerified: updatedAuthor.IsVerified,
		CreatedAt:  updatedAuthor.CreatedAt,
		UpdatedAt:  updatedAuthor.UpdatedAt,
	}, nil
}

func (s *AuthorService) ListAuthors(ctx interface{}, limit, offset int) ([]*models.AuthorResponse, error) {
	authors, err := s.store.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.AuthorResponse, len(authors))
	for i, author := range authors {
		responses[i] = &models.AuthorResponse{
			ID:         author.ID,
			Username:   author.Username,
			Email:      author.Email,
			FirstName:  author.FirstName,
			LastName:   author.LastName,
			Bio:        author.Bio,
			AvatarURL:  author.AvatarURL,
			IsActive:   author.IsActive,
			IsVerified: author.IsVerified,
			CreatedAt:  author.CreatedAt,
			UpdatedAt:  author.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *AuthorService) DeleteAccount(ctx interface{}, authorID int) error {
	return s.store.Delete(ctx, authorID)
}

func (s *AuthorService) ValidateToken(tokenString string) (*models.Author, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract author information from claims
	authorID, ok := claims["author_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &models.Author{
		ID:       int(authorID),
		Username: username,
	}, nil
}

func (s *AuthorService) generateJWT(author *models.Author) (string, error) {
	claims := jwt.MapClaims{
		"author_id": author.ID,
		"username":  author.Username,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
