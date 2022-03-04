package app

import (
	"cro_test/internal/delivery/http"
	"cro_test/internal/repository"
	"cro_test/internal/repository/mysql"
	"cro_test/internal/service"
	"cro_test/pkg/config"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewServerApp(cfgName string, opts ...fx.Option) (*fx.App, error) {
	err := config.Init(cfgName)
	if err != nil {
		return nil, err
	}
	if opts == nil {
		opts = make([]fx.Option, 0, 1)
	}
	opts = append(opts, fx.Options(
		fx.Provide(config.InitLogger),
		ProvideSvc(),
		fx.Provide(http.New),
	))
	if viper.Get("app.env") == "prod" {
		opts = append(opts, fx.NopLogger)
	}
	app := fx.New(opts...)
	return app, nil
}

func ProvideSvc() fx.Option {
	return fx.Options(
		ProvideRepo(),
		fx.Provide(service.New),
	)
}

func ProvideRepo() fx.Option {
	return fx.Options(
		fx.Provide(config.InitSqlxx),
		fx.Provide(mysql.NewDao),
		fx.Provide(repository.New),
	)
}
