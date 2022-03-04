package model

import "github.com/shopspring/decimal"

type TxKind uint8

const (
	TxKindDeposit TxKind = iota + 1
	TxKindWithdraw
	TxKindTransfer
)

type Transaction struct {
	ID                uint64          `db:"id" json:"id"`
	Kind              TxKind          `db:"kind" json:"kind"`
	FromWalletID      uint64          `db:"from_wallet_id" json:"fromWalletID"`
	ToWalletID        uint64          `db:"to_wallet_id" json:"toWalletID"`
	FromWalletBalance decimal.Decimal `json:"-"`
	ToWalletBalance   decimal.Decimal `json:"-"`
	TxAmount          decimal.Decimal `db:"tx_amount" json:"txAmount"`
	CreatedAti        int64           `db:"created_ati" json:"createdAti"`
}

func (tx *Transaction) SetBalance(from, to Wallet) {
	switch tx.Kind {
	case TxKindTransfer:
		tx.FromWalletBalance = from.Amount
		tx.ToWalletBalance = to.Amount
	case TxKindDeposit:
		tx.ToWalletBalance = to.Amount
	case TxKindWithdraw:
		tx.FromWalletBalance = from.Amount
	}
}
