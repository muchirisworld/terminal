package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DBTX defines the common methods for both *sqlx.DB and *sqlx.Tx.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Store is the main repository that provides database operations.
type Store struct {
	db   *sqlx.DB
	dbtx DBTX
}

// New creates a new Store.
func New(db *sqlx.DB) *Store {
	return &Store{
		db:   db,
		dbtx: db,
	}
}

// WithTx returns a new Store instance that uses the provided transaction.
func (s *Store) WithTx(tx *sqlx.Tx) *Store {
	return &Store{
		db:   s.db,
		dbtx: tx,
	}
}

// ExecTx executes a function within a database transaction.
func (s *Store) ExecTx(ctx context.Context, fn func(*Store) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	q := s.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
