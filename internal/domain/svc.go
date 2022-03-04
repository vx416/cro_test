package domain

import (
	"context"
	"database/sql"
)

type Servicer interface {
	WalletServicer
	TransactionServicer
	AuthServicer
}

type Repositorier interface {
	TransactionRepositorier
	WalletRepositorier
	AuthRepositorier
	ExecuteTx(ctx context.Context, fn func(txCtx context.Context) error, txOpts ...*sql.TxOptions) error
}
