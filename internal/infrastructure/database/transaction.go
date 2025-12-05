package database

import (
	"context"

	"gorm.io/gorm"
)

type TransactionSession interface {
	Commit() error
	Rollback() error
}

type GormTransaction struct {
	Tx *gorm.DB
}

func (t *GormTransaction) Commit() error {
	return t.Tx.Commit().Error
}

func (t *GormTransaction) Rollback() error {
	return t.Tx.Rollback().Error
}

// Transaction Runner
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
