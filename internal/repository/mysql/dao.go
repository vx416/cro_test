package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/vx416/sqlxx"
)

var (
	ErrNotAffected = errors.New("rows not affected")
)

func NewDao(sqlxAdapter *sqlxx.Sqlxx) MySQLDao {
	return MySQLDao{
		Sqlxx: sqlxAdapter,
	}
}

type MySQLDao struct {
	*sqlxx.Sqlxx
}

func IsDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)

	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		if mysqlErr.Number == 1062 {
			return true
		}
	}

	return false
}
