package postgres

import (
	"context"
	"database/sql"
)

// Driver is postgres database driver implementation
type Driver struct {
	db *sql.DB
}

// NewDriver is Driver constructor
func NewDriver(ctx context.Context, db *sql.DB) *Driver {
	return &Driver{
		db: db,
	}
}
