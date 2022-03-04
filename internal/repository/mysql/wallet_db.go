package mysql

import (
	"context"
	"cro_test/internal/model"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vx416/sqlxx/builder"
)

func (dao MySQLDao) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	db := dao.GetDB(ctx)
	if wallet.CreatedAti == 0 || wallet.UpdatedAti == 0 {
		wallet.CreatedAti = time.Now().Unix()
		wallet.UpdatedAti = time.Now().Unix()
	}
	insertStmt := builder.Insert()
	insertStmt.InsertRows(wallet).Table("wallets")
	res, err := db.Exec(ctx, insertStmt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	wallet.ID = uint64(id)
	return nil
}

func (dao MySQLDao) GetWalletBySerial(ctx context.Context, serialNumber string) (model.Wallet, error) {
	wallet := model.Wallet{}
	db := dao.GetDB(ctx)
	query := builder.Query().And("serial_number = ?", serialNumber).From("wallets")
	err := db.Get(ctx, &wallet, query)
	if err != nil {
		return model.Wallet{}, err
	}

	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	txs := []model.Transaction{}
	txQuery := builder.Query().From("transactions").And("from_wallet_id = ?", wallet.ID).Or("to_wallet_id = ?", wallet.ID).
		And("created_ati >= ?", weekAgo.Unix())
	err = db.Select(ctx, &txs, txQuery)
	if err != nil {
		return model.Wallet{}, err
	}
	wallet.Transactions = txs
	return wallet, nil
}

func (dao MySQLDao) ListWallets(ctx context.Context, opts model.ListWalletsOpts) ([]model.Wallet, error) {
	wallets := make([]model.Wallet, 0, 10)
	db := dao.GetDB(ctx)
	query := builder.Query().Where(opts, builder.SkipZero).From("wallets")

	err := db.Select(ctx, &wallets, query)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// CreateWallet(ctx context.Context, wallet model.Wallet) error
func (dao MySQLDao) UpdateWalletAmountDelta(ctx context.Context, serialNumber string, amountDelta decimal.Decimal) (model.Wallet, error) {
	db := dao.GetDB(ctx)
	wallet := model.Wallet{}
	updateStmt := builder.Update().Set("updated_ati = ?", time.Now().Unix()).And("serial_number = ?", serialNumber).Table("wallets")
	if amountDelta.IsNegative() {
		amountDelta = amountDelta.Abs()
		updateStmt.And("amount >= ?", amountDelta).Set("amount = amount - ?", amountDelta)
	} else {
		updateStmt.Set("amount = amount + ?", amountDelta)
	}

	res, err := db.Exec(ctx, updateStmt)
	if err != nil {
		return model.Wallet{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return model.Wallet{}, err
	}
	if rowsAffected == 0 {
		return wallet, ErrNotAffected
	}

	return dao.GetWalletBySerial(ctx, serialNumber)
}
