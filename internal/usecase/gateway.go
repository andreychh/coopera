package usecase

import (
	"context"
	"github.com/jackc/pgx/v4"
)

// Здесь будут заданы интерфейсы для взаимодействия с имплементациями сдоя репозитория и управления транзакциями

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	WithinTransactionWithIsolation(ctx context.Context, level pgx.TxIsoLevel, fn func(txCtx context.Context) error) error
}
