package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrMissingParameter = errors.New("missing parameter")
	ErrNotFound         = errors.New("record not found")
)

// Driver is postgres database driver implementation
type Driver struct {
	db *sql.DB
	l  *zap.Logger
}

// NewDriver is Driver constructor
func NewDriver(ctx context.Context, db *sql.DB, l *zap.Logger) *Driver {
	return &Driver{
		db: db,
		l:  l,
	}
}
