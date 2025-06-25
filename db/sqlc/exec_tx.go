package db

import (
	"context"
	"fmt"
)

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Start a new transaction
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		// Rollback the transaction if there was an error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit the transaction
	return tx.Commit(ctx)
}
