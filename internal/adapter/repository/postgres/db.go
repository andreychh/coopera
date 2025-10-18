package postgres

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera/internal/usecase"
	"github.com/andreychh/coopera/pkg/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ usecase.TransactionManager = (*DB)(nil)

type DB struct {
	Pool   *pgxpool.Pool
	Logger *logger.Logger
}

type TransactionKey struct{}

func NewDB(dsn string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DB{Pool: pool}, nil
}

func (db *DB) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return db.WithinTransactionWithIsolation(ctx, pgx.ReadCommitted, fn)
}

func (db *DB) WithinTransactionWithIsolation(ctx context.Context, level pgx.TxIsoLevel, fn func(txCtx context.Context) error) error {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: level})
	if err != nil {
		return err
	}

	wrappedTx := &PGXTransaction{tx: tx}
	ctx = context.WithValue(ctx, TransactionKey{}, wrappedTx)

	defer func() {
		if p := recover(); p != nil {
			_ = wrappedTx.Rollback(ctx)
			err = fmt.Errorf("transaction failed with panic and rollback: %v", p)
		} else if err != nil {
			_ = wrappedTx.Rollback(ctx)
		} else {
			err = wrappedTx.Commit(ctx)
		}
	}()

	err = fn(ctx)
	return err
}

func (db *DB) Close() {
	db.Pool.Close()
}
