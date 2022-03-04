package service

import (
	"context"
	"cro_test/internal/model"
	"cro_test/internal/repository/mysql"
	"cro_test/pkg/httperr"
	"database/sql"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"
)

func (svc Service) Transfer(ctx context.Context, userID uint64, fromWalletNumber, toWalletNumber string, amount decimal.Decimal) (model.Transaction, error) {
	var (
		fromWallet = model.Wallet{}
		toWallet   = model.Wallet{}
		tx         = model.Transaction{}
		err        error
	)

	err = svc.repo().ExecuteTx(ctx, func(txCtx context.Context) error {
		wallets, err := svc.repo().ListWallets(txCtx, model.ListWalletsOpts{
			SerilaNumbers: []string{fromWalletNumber, toWalletNumber},
		})
		if err != nil {
			return err
		}
		if len(wallets) != 2 {
			return httperr.WrapfErr(http.StatusUnprocessableEntity, "serialize number not found", "from:%s, to:%s", fromWalletNumber, toWalletNumber)
		}
		if wallets[0].SerialNumber == fromWalletNumber {
			fromWallet = wallets[0]
			toWallet = wallets[1]
		} else {
			fromWallet = wallets[1]
			toWallet = wallets[0]
		}
		if !fromWallet.CanUse(userID) {
			return httperr.WrapfErr(http.StatusUnauthorized, "not wallet owner", "actual:%d, expect:%d", fromWallet.UserID, userID)
		}
		tx, err = fromWallet.TransferTo(toWallet, amount)
		if err != nil {
			return httperr.WrapfErr(http.StatusUnprocessableEntity, "currency invalid", "err:%+v", err)
		}

		err = svc.repo().CreateTransaction(txCtx, tx)
		if err != nil {
			return httperr.WrapfErr(http.StatusInternalServerError, "", "create tx failed, err:%+v", err)
		}
		fromWallet, err = svc.repo().UpdateWalletAmountDelta(txCtx, fromWalletNumber, amount.Neg())
		if err != nil {
			if errors.Is(err, mysql.ErrNotAffected) {
				return httperr.WrapfErr(http.StatusUnprocessableEntity, "wallet balance insufficient", "current:%s, delta:%s", fromWallet.Amount, amount)
			}
			return err
		}
		toWallet, err = svc.repo().UpdateWalletAmountDelta(txCtx, toWalletNumber, amount)
		if err != nil {
			return httperr.WrapfErr(http.StatusInternalServerError, "", "increase toWallet amount failed, err:%+v", err)
		}
		tx.SetBalance(fromWallet, toWallet)
		return nil
	})
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (svc Service) Deposit(ctx context.Context, userID uint64, toWalletNumber string, amount decimal.Decimal) (model.Transaction, error) {
	if amount.IsNegative() {
		amount = amount.Abs()
	}

	return svc.checkAndUpdateWalletAmount(ctx, userID, toWalletNumber, amount, model.TxKindDeposit)
}

func (svc Service) Withdraw(ctx context.Context, userID uint64, fromWalletNumber string, amount decimal.Decimal) (model.Transaction, error) {
	if amount.IsPositive() {
		amount = amount.Neg()
	}
	return svc.checkAndUpdateWalletAmount(ctx, userID, fromWalletNumber, amount, model.TxKindWithdraw)
}

func (svc Service) checkAndUpdateWalletAmount(ctx context.Context, userID uint64, walletNumber string, amount decimal.Decimal, txKind model.TxKind) (model.Transaction, error) {
	var (
		wallet = model.Wallet{}
		tx     = model.Transaction{}
		err    error
	)

	err = svc.repo().ExecuteTx(ctx, func(txCtx context.Context) error {
		wallet, err = svc.repo().GetWalletBySerial(txCtx, walletNumber)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return httperr.WrapfErr(http.StatusUnprocessableEntity, "serialize number not found", "number:%s", walletNumber)
			}
			return err
		}
		if !wallet.CanUse(userID) {
			return httperr.WrapfErr(http.StatusUnauthorized, "not wallet owner", "actual:%d, expect:%d", wallet.UserID, userID)
		}
		tx, err = wallet.DepositOrWithdraw(txKind, amount)
		if err != nil {
			return httperr.WrapfErr(http.StatusInternalServerError, "cannot get transaction", "err:%+v", err)
		}
		err = svc.repo().CreateTransaction(txCtx, tx)
		if err != nil {
			return httperr.WrapfErr(http.StatusInternalServerError, "", "create tx failed, err:%+v", err)
		}
		wallet, err = svc.repo().UpdateWalletAmountDelta(txCtx, walletNumber, amount)
		if err != nil {
			if errors.Is(err, mysql.ErrNotAffected) {
				return httperr.WrapfErr(http.StatusUnprocessableEntity, "wallet balance insufficient", "current:%s, delta:%s", wallet.Amount, amount)
			}
			return err
		}
		tx.SetBalance(wallet, wallet)
		return nil
	})
	if err != nil {
		return tx, err
	}
	return tx, nil
}
