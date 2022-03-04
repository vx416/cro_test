package service

import (
	"context"
	"cro_test/internal/model"
	"cro_test/internal/repository/mysql"
	"cro_test/pkg/httperr"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (svc Service) CreateWallet(ctx context.Context, user *model.User, currency model.CurrencyValue) (model.Wallet, error) {
	wallet := model.Wallet{}
	err := svc.repo().ExecuteTx(ctx, func(txCtx context.Context) error {
		err := svc.findUserByID(txCtx, user)
		if err != nil {
			return err
		}
		wallet = model.Wallet{
			UserID:       user.ID,
			SerialNumber: uuid.New().String(),
			Amount:       decimal.Zero,
			Currency:     currency,
		}
		err = svc.repo().CreateWallet(txCtx, &wallet)
		if err != nil {
			if mysql.IsDuplicateErr(err) {
				return httperr.WrapfErr(http.StatusConflict, "currency duplicated", "currency:%s", currency)
			}
			return httperr.WrapfErr(http.StatusInternalServerError, "", "create wallet failed, err:%+v", err)
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return wallet, nil
}

func (svc Service) GetWalletBySerial(ctx context.Context, userID uint64, serialNumber string) (model.Wallet, error) {
	wallet, err := svc.repo().GetWalletBySerial(ctx, serialNumber)
	if err != nil {
		return model.Wallet{}, err
	}
	if !wallet.CanUse(userID) {
		return wallet, httperr.WrapfErr(http.StatusUnauthorized, "not wallet owner", "actual:%d, expect:%d", wallet.UserID, userID)
	}
	return wallet, nil
}

func (svc Service) findUserByID(ctx context.Context, user *model.User) error {
	currUser, err := svc.repo().GetUser(ctx, model.GetUserOpts{
		ID: user.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return httperr.WrapfErr(http.StatusUnauthorized, "use not found", "id:%d", user.ID)
		}
		return httperr.WrapfErr(http.StatusInternalServerError, "", "get user failed, err:%+v", err)
	}
	user.Email = currUser.Email
	user.ID = currUser.ID
	user.CreatedAti = currUser.CreatedAti
	return nil
}
