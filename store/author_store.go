package store

import (
	"database/sql"
	"errors"

	"gofr.dev/pkg/gofr"
	"gofr-blog-service/models"
)

var (
	errAuthorNotFound = errors.New("author not found")
	errUsernameExists = errors.New("username already exists")
	errEmailExists    = errors.New("email already exists")
)

type AuthorStore struct{}

func NewAuthorStore() *AuthorStore {
	return &AuthorStore{}
}

func (s *AuthorStore) Create(ctx interface{}, author *models.Author) (*models.Author, error) {
	gofrCtx := ctx.(*gofr.Context)
	var id int
	var createdAuthor models.Author

	err := gofrCtx.SQL.QueryRow(createAuthorQuery,
		author.Username, author.Email, author.PasswordHash,
		author.FirstName, author.LastName, author.Bio, author.AvatarURL).
		Scan(&id, &createdAuthor.Username, &createdAuthor.Email,
			&createdAuthor.FirstName, &createdAuthor.LastName,
			&createdAuthor.Bio, &createdAuthor.AvatarURL,
			&createdAuthor.IsActive, &createdAuthor.IsVerified,
			&createdAuthor.CreatedAt, &createdAuthor.UpdatedAt)

	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "authors_username_key"` {
			return nil, errUsernameExists
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "authors_email_key"` {
			return nil, errEmailExists
		}
		return nil, err
	}

	createdAuthor.ID = id
	return &createdAuthor, nil
}

func (s *AuthorStore) GetByID(ctx interface{}, id int) (*models.Author, error) {
	gofrCtx := ctx.(*gofr.Context)
	var author models.Author

	err := gofrCtx.SQL.QueryRow(getAuthorByIDQuery, id).
		Scan(&author.ID, &author.Username, &author.Email,
			&author.FirstName, &author.LastName, &author.Bio,
			&author.AvatarURL, &author.IsActive, &author.IsVerified,
			&author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errAuthorNotFound
		}
		return nil, err
	}

	return &author, nil
}

func (s *AuthorStore) GetByUsername(ctx interface{}, username string) (*models.Author, error) {
	gofrCtx := ctx.(*gofr.Context)
	var author models.Author

	err := gofrCtx.SQL.QueryRow(getAuthorByUsernameQuery, username).
		Scan(&author.ID, &author.Username, &author.Email,
			&author.FirstName, &author.LastName, &author.Bio,
			&author.AvatarURL, &author.IsActive, &author.IsVerified,
			&author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errAuthorNotFound
		}
		return nil, err
	}

	return &author, nil
}

func (s *AuthorStore) GetByEmail(ctx interface{}, email string) (*models.Author, error) {
	gofrCtx := ctx.(*gofr.Context)
	var author models.Author

	err := gofrCtx.SQL.QueryRow(getAuthorByEmailQuery, email).
		Scan(&author.ID, &author.Username, &author.Email,
			&author.FirstName, &author.LastName, &author.Bio,
			&author.AvatarURL, &author.IsActive, &author.IsVerified,
			&author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errAuthorNotFound
		}
		return nil, err
	}

	return &author, nil
}

func (s *AuthorStore) GetPasswordHash(ctx interface{}, username string) (string, error) {
	gofrCtx := ctx.(*gofr.Context)
	var passwordHash string

	err := gofrCtx.SQL.QueryRow(getPasswordHashQuery, username).
		Scan(&passwordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errAuthorNotFound
		}
		return "", err
	}

	return passwordHash, nil
}

func (s *AuthorStore) Update(ctx interface{}, id int, author *models.Author) (*models.Author, error) {
	gofrCtx := ctx.(*gofr.Context)
	var updatedAuthor models.Author

	err := gofrCtx.SQL.QueryRow(updateAuthorQuery,
		author.Email, author.FirstName, author.LastName,
		author.Bio, author.AvatarURL, id).
		Scan(&updatedAuthor.ID, &updatedAuthor.Username, &updatedAuthor.Email,
			&updatedAuthor.FirstName, &updatedAuthor.LastName,
			&updatedAuthor.Bio, &updatedAuthor.AvatarURL,
			&updatedAuthor.IsActive, &updatedAuthor.IsVerified,
			&updatedAuthor.CreatedAt, &updatedAuthor.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errAuthorNotFound
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "authors_email_key"` {
			return nil, errEmailExists
		}
		return nil, err
	}

	return &updatedAuthor, nil
}

func (s *AuthorStore) Delete(ctx interface{}, id int) error {
	gofrCtx := ctx.(*gofr.Context)
	result, err := gofrCtx.SQL.Exec(deleteAuthorQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errAuthorNotFound
	}

	return nil
}

func (s *AuthorStore) List(ctx interface{}, limit, offset int) ([]*models.Author, error) {
	gofrCtx := ctx.(*gofr.Context)
	rows, err := gofrCtx.SQL.Query(listAuthorsQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*models.Author

	for rows.Next() {
		var author models.Author
		err := rows.Scan(&author.ID, &author.Username, &author.Email,
			&author.FirstName, &author.LastName, &author.Bio,
			&author.AvatarURL, &author.IsActive, &author.IsVerified,
			&author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, &author)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}
