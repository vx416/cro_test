package domain

import (
	"context"
	"cro_test/internal/model"

	"github.com/shopspring/decimal"
)

type WalletRepositorier interface {
	GetWalletBySerial(ctx context.Context, serialNumber string) (model.Wallet, error)
	ListWallets(ctx context.Context, opts model.ListWalletsOpts) ([]model.Wallet, error)
	CreateWallet(ctx context.Context, wallet *model.Wallet) error
	UpdateWalletAmountDelta(ctx context.Context, serialNumber string, amountDelta decimal.Decimal) (model.Wallet, error)
}

type WalletServicer interface {
	ListWallets(ctx context.Context, opts model.ListWalletsOpts) ([]model.Wallet, error)
	GetWalletBySerial(ctx context.Context, userID uint64, serialNumber string) (model.Wallet, error)
	CreateWallet(ctx context.Context, user *model.User, currency model.CurrencyValue) (model.Wallet, error)
}
