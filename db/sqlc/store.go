package db

import (
	"database/sql"
)

// Store defines all functions to execute db queries and transactions
type Store struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}
