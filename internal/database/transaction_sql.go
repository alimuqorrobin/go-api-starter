package transaction

import (
	"context"
	"database/sql"
)

// TxFuncSQL signature
type TxFuncSQL func(tx *sql.Tx) error

func WithTransactionSQL(ctx context.Context, db *sql.DB, fn TxFuncSQL) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
