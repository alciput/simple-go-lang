package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"app/models"
	"app/store"

	_ "modernc.org/sqlite"
)

// SQLite implements store.UserRepository using SQLite.
type SQLite struct {
	db *sql.DB
}

// New connects to SQLite at path and creates the users table if needed.
func New(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	s := &SQLite{db: db}
	if err := s.createTable(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *SQLite) createTable(ctx context.Context) error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.ExecContext(ctx, stmt)
	return err
}

// CreateUser inserts a user. Returns store.ErrDuplicateEmail if email already exists.
func (s *SQLite) CreateUser(ctx context.Context, email, passwordHash string) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash) VALUES (?, ?)`,
		email, passwordHash,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, store.ErrDuplicateEmail
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// GetUserByEmail returns the user or nil if not found.
func (s *SQLite) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// Close closes the database connection.
func (s *SQLite) Close() error {
	return s.db.Close()
}

// Ensure SQLite implements store.UserRepository.
var _ store.UserRepository = (*SQLite)(nil)
