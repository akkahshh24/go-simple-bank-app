package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute SQL queries and transaction
type Store struct {
	*Queries // compose the queries struct to extend the functionality
	connPool *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  New(connPool),
		connPool: connPool,
	}
}
