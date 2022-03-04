package mysql

import (
	"context"
	"cro_test/internal/model"
	"time"

	"github.com/vx416/sqlxx/builder"
)

func (dao MySQLDao) GetUser(ctx context.Context, opts model.GetUserOpts) (model.User, error) {
	db := dao.GetDB(ctx)
	user := model.User{}

	err := db.Get(ctx, &user, builder.Query().From("users").Where(opts, builder.SkipZero))
	if err != nil {
		return user, err
	}
	return user, nil
}

func (dao MySQLDao) CreateUser(ctx context.Context, user *model.User) error {
	db := dao.GetDB(ctx)
	if user.CreatedAti == 0 {
		user.CreatedAti = time.Now().Unix()
	}
	insertStmt := builder.Insert()
	insertStmt.InsertRows(user).Table("users")
	res, err := db.Exec(ctx, insertStmt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(id)
	return nil
}
