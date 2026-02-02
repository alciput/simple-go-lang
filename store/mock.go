package store

import (
	"context"
	"sync"
	"time"

	"app/models"
)

// Mock is an in-memory UserRepository for tests.
type Mock struct {
	mu    sync.Mutex
	users map[string]*models.User
	nextID int64
}

// NewMock returns a Mock with empty user map.
func NewMock() *Mock {
	return &Mock{users: make(map[string]*models.User), nextID: 1}
}

// CreateUser stores a user by email. Returns ErrDuplicateEmail if email exists.
func (m *Mock) CreateUser(ctx context.Context, email, passwordHash string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[email]; ok {
		return 0, ErrDuplicateEmail
	}
	id := m.nextID
	m.nextID++
	m.users[email] = &models.User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now()}
	return id, nil
}

// GetUserByEmail returns the user or nil if not found.
func (m *Mock) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[email]
	if !ok {
		return nil, nil
	}
	// Return a copy so callers cannot mutate stored data.
	cp := *u
	return &cp, nil
}
