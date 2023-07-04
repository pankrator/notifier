package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pankrator/notifier/internal/storage"
)

type txKeyType string

const txKey = txKeyType("db.tx")

type Transactioner interface {
	storage.Queryer
	RunInTx(context.Context, func(ctx context.Context) error) error
}

func ExtractTransactionFromContext(ctx context.Context) *Tx {
	if tx, ok := ctx.Value(txKey).(*Tx); ok {
		return tx
	}

	return nil
}

func putTransactionToContext(
	ctx context.Context,
	tx *Tx,
) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

type Tx struct {
	*sqlx.Tx
}
