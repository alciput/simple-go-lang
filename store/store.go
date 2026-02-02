package store

import (
	"context"
	"errors"

	"app/models"
)

// ErrDuplicateEmail is returned when CreateUser is called with an existing email.
var ErrDuplicateEmail = errors.New("email already registered")

// UserRepository abstracts user persistence for handlers and tests.
type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
