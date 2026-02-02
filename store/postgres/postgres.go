package postgres

import (
	"context"
	"errors"
	"fmt"

	"app/models"
	"app/store"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres implements store.UserRepository using pgxpool.
type Postgres struct {
	pool *pgxpool.Pool
}

// New connects to PostgreSQL and creates the users table if needed.
func New(ctx context.Context, databaseURL string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool new: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}
	p := &Postgres{pool: pool}
	if err := p.createTable(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return p, nil
}

func (p *Postgres) createTable(ctx context.Context) error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)`
	_, err := p.pool.Exec(ctx, stmt)
	return err
}

// CreateUser inserts a user. Returns store.ErrDuplicateEmail if email already exists.
func (p *Postgres) CreateUser(ctx context.Context, email, passwordHash string) (int64, error) {
	var id int64
	err := p.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		email, passwordHash,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, store.ErrDuplicateEmail
		}
		return 0, err
	}
	return id, nil
}

// GetUserByEmail returns the user or nil if not found.
func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := p.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// Close closes the connection pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// Ensure Postgres implements store.UserRepository.
var _ store.UserRepository = (*Postgres)(nil)
