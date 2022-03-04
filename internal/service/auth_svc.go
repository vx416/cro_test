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
)

func (svc Service) SignUp(ctx context.Context, email, password string) (model.User, error) {
	user := model.User{
		Email:    email,
		Password: password,
	}
	err := user.HashPwd()
	if err != nil {
		return model.User{}, httperr.WrapfErr(http.StatusInternalServerError, "", "hash password failed, err:%+v", err)
	}

	err = svc.repo().ExecuteTx(ctx, func(txCtx context.Context) error {
		err = svc.repo().CreateUser(ctx, &user)
		if err != nil {
			if mysql.IsDuplicateErr(err) {
				return httperr.WrapfErr(http.StatusConflict, "email duplicated", "email:%s", email)
			}
			return err
		}
		svc.repo().CreateWallet(txCtx, &model.Wallet{
			UserID:       user.ID,
			SerialNumber: uuid.New().String(),
		})
		return nil
	})
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (svc Service) Login(ctx context.Context, email, password string) (model.User, error) {
	currUser, err := svc.GetUser(ctx, model.GetUserOpts{
		Email: email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, httperr.WrapfErr(http.StatusBadRequest, "email not found", "email:%s", email)
		}
		return model.User{}, err
	}
	err = currUser.CmpPwd(password)
	if err != nil {
		return model.User{}, httperr.WrapfErr(http.StatusUnauthorized, "password incorrect", "err:%+v", err)
	}
	return currUser, nil
}
