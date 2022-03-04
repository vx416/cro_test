package mysql

import (
	"context"
	"cro_test/internal/model"
	"time"

	"github.com/vx416/sqlxx/builder"
)

func (dao MySQLDao) CreateTransaction(ctx context.Context, tx model.Transaction) error {
	if tx.CreatedAti == 0 {
		tx.CreatedAti = time.Now().Unix()
	}
	db := dao.GetDB(ctx)
	insertStmt := builder.Insert()
	insertStmt.InsertRows(tx).Table("transactions")
	_, err := db.Exec(ctx, insertStmt)
	if err != nil {
		return err
	}
	return nil
}
