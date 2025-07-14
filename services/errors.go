package services

import "errors"

// Error definitions for the service layer
var (
	ErrCreateFailed     = errors.New("failed to create post")
	ErrGetFailed        = errors.New("failed to get post")
	ErrListFailed       = errors.New("failed to list posts")
	ErrCountFailed      = errors.New("failed to count posts")
	ErrUpdateFailed     = errors.New("failed to update post")
	ErrDeleteFailed     = errors.New("failed to delete post")
	ErrValidationFailed = errors.New("validation failed")
)
