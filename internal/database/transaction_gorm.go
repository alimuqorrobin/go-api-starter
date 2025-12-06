package transaction

import (
	"context"

	"gorm.io/gorm"
)

// TxFuncGORM signature
type TxFuncGORM func(tx *gorm.DB) error

func WithTransactionGORM(ctx context.Context, db *gorm.DB, fn TxFuncGORM) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
