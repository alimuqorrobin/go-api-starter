package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause" // ← dari gorm, bukan terpisah
)

// TransactionFunc adalah function yang dijalankan dalam transaction
type TransactionFunc func(*gorm.DB) error

// WithTx wraps function dalam transaction
func WithTx(db *gorm.DB, ctx context.Context, fn TransactionFunc) error {
	if db == nil {
		return fmt.Errorf("database instance is nil")
	}
	if fn == nil {
		return fmt.Errorf("transaction function is nil")
	}

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTxRecovery wraps function dengan panic recovery
func WithTxRecovery(db *gorm.DB, ctx context.Context, fn TransactionFunc) (err error) {
	if db == nil {
		return fmt.Errorf("database instance is nil")
	}
	if fn == nil {
		return fmt.Errorf("transaction function is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in transaction: %v", r)
		}
	}()

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ChainTx chains multiple transactions
func ChainTx(db *gorm.DB, ctx context.Context, fns ...TransactionFunc) error {
	return WithTx(db, ctx, func(tx *gorm.DB) error {
		for i, fn := range fns {
			if err := fn(tx); err != nil {
				return fmt.Errorf("transaction step %d failed: %w", i+1, err)
			}
		}
		return nil
	})
}

// ============================================
// LOCKING HELPERS - Race Condition Protection
// ============================================

// LockForUpdate locks row untuk update (pessimistic locking)
func LockForUpdate(tx *gorm.DB, model interface{}, id uint) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(model, id).Error
}

// LockForUpdateWhere locks row dengan custom where clause
func LockForUpdateWhere(tx *gorm.DB, model interface{}, where string, args ...interface{}) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(where, args...).
		First(model).Error
}

// LockForShare locks row untuk read (shared lock)
func LockForShare(tx *gorm.DB, model interface{}, id uint) error {
	return tx.Clauses(clause.Locking{Strength: "SHARE"}).
		First(model, id).Error
}

// LockForShareWhere locks row dengan custom where clause (shared)
func LockForShareWhere(tx *gorm.DB, model interface{}, where string, args ...interface{}) error {
	return tx.Clauses(clause.Locking{Strength: "SHARE"}).
		Where(where, args...).
		First(model).Error
}

// NoWait tries lock, fail immediately if locked
func NoWait(tx *gorm.DB, model interface{}, id uint) error {
	return tx.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "NOWAIT",
	}).First(model, id).Error
}

// SkipLocked skips locked rows
func SkipLocked(tx *gorm.DB, model interface{}) *gorm.DB {
	return tx.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	})
}