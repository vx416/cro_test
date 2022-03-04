package domain

import (
	"context"
	"cro_test/internal/model"

	"github.com/shopspring/decimal"
)

type TransactionRepositorier interface {
	CreateTransaction(context.Context, model.Transaction) error
}

type TransactionServicer interface {
	Transfer(ctx context.Context, userID uint64, fromWalletNumber, toWalletNumber string, amount decimal.Decimal) (model.Transaction, error)
	Deposit(ctx context.Context, userID uint64, toWalletNumber string, amount decimal.Decimal) (model.Transaction, error)
	Withdraw(ctx context.Context, userID uint64, fromWalletNumber string, amount decimal.Decimal) (model.Transaction, error)
}
