package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB is the shared database handle. Use prepared statements only; never concatenate user input.
var DB *sql.DB

// Connect opens SQLite at path and creates tables if they don't exist.
func Connect(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	DB = db
	return createTables()
}

// createTables runs idempotent DDL. Uses parameterless statements; table names are fixed.
func createTables() error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := DB.Exec(stmt)
	return err
}

// Close closes the database connection.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
