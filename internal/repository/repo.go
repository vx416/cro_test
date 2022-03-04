package repository

import (
	"cro_test/internal/domain"
	"cro_test/internal/repository/mysql"

	"go.uber.org/fx"
)

type Params struct {
	fx.In
	MySQLDao mysql.MySQLDao
}

func New(params Params) domain.Repositorier {
	return Repository{
		MySQLDao: params.MySQLDao,
	}
}

type Repository struct {
	mysql.MySQLDao
}
