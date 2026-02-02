package models

import "time"

// User represents a row in the users table. password_hash is never returned in JSON.
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never expose in API
	CreatedAt    time.Time `json:"created_at"`
}
