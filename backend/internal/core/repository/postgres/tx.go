package postgres_pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryExecutor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type txKeyType struct{}

var txKey txKeyType

type TxManager struct {
	pool Pool
}

func NewTxManager(pool Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (tm *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.do(ctx, fn, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
}

func (tm *TxManager) DoWrite(ctx context.Context, fn func(context.Context) error) error {
	return tm.do(ctx, fn, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
}

func (tm *TxManager) do(ctx context.Context, fn func(context.Context) error, options pgx.TxOptions) error {
	if _, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := tm.pool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return fmt.Errorf("rollback err: %w: %w", rollbackErr, err)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func ExecutorFromContext(ctx context.Context, fallback QueryExecutor) QueryExecutor {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return fallback
}
