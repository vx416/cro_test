package model

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v3"
)

type CurrencyValue string

func (val CurrencyValue) Eq(s CurrencyValue) bool {
	return strings.EqualFold(string(val), string(s))
}

type ListWalletsOpts struct {
	UserID        uint64   `sql:"col:user_id"`
	SerilaNumbers []string `sql:"col:serial_number;op:in"`
}

type Wallet struct {
	ID           uint64          `db:"id" json:"-"`
	SerialNumber string          `db:"serial_number" json:"serialNumber"`
	Currency     CurrencyValue   `db:"currency" json:"currency"`
	UserID       uint64          `db:"user_id" json:"userID"`
	Amount       decimal.Decimal `db:"amount" json:"amount"`
	CreatedAti   int64           `db:"created_ati" json:"createdAti"`
	UpdatedAti   int64           `db:"updated_ati" json:"updatedAti"`
	DeltedAti    null.Int        `db:"deleted_ati" json:"-"`
	Transactions []Transaction   `json:"transactions,omitempty"`
}

func (wallet Wallet) CanUse(userID uint64) bool {
	return wallet.UserID == userID
}

func (wallet Wallet) TransferTo(toWallet Wallet, amount decimal.Decimal) (Transaction, error) {
	if !wallet.Currency.Eq(toWallet.Currency) {
		return Transaction{}, fmt.Errorf("currency is not the same, from:%s, to:%s", wallet.Currency, toWallet.Currency)
	}
	return Transaction{
		FromWalletID: wallet.ID,
		ToWalletID:   toWallet.ID,
		TxAmount:     amount,
		Kind:         TxKindTransfer,
	}, nil
}

func (wallet Wallet) DepositOrWithdraw(kind TxKind, amount decimal.Decimal) (Transaction, error) {
	tx := Transaction{}

	switch kind {
	case TxKindDeposit:
		return Transaction{
			ToWalletID: wallet.ID,
			TxAmount:   amount.Abs(),
			Kind:       TxKindDeposit,
		}, nil
	case TxKindWithdraw:
		return Transaction{
			FromWalletID: wallet.ID,
			TxAmount:     amount.Abs(),
			Kind:         TxKindWithdraw,
		}, nil
	}

	return tx, fmt.Errorf("unknown kind(%d)", kind)
}
