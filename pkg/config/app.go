package config

import (
	"context"
	"cro_test/pkg/httperr"
	"cro_test/pkg/logger"
	h "net/http"
	"path/filepath"
	"runtime"

	echoSwagger "github.com/swaggo/echo-swagger"

	mid "cro_test/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"github.com/vx416/sqlxx"
	"go.uber.org/fx"
)

func RunMigrations(db *sqlxx.Sqlxx) error {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	migPath := filepath.Join(dir, "../../build/migrations")
	sqlDB, err := db.GetDB(context.Background()).GetRawDB(context.Background())
	if err != nil {
		return err
	}
	goose.SetDialect("mysql")
	return goose.Run("up", sqlDB, migPath)
}

func StartEchoServer(lc fx.Lifecycle, log logger.Logger) (*echo.Echo, error) {
	e := echo.New()
	reqDump := false
	respDump := false
	if viper.GetString("app.env") == "dev" {
		e.Debug = true
		e.HideBanner = false
		e.HidePort = false
		reqDump = true
		respDump = true
	} else {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	}

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(h.StatusNotFound, map[string]interface{}{"error": "Page Not Found"})
	}

	setupSwagger(e)
	e.HTTPErrorHandler = httperr.EchoErrHandle
	e.Use(echo.WrapMiddleware(mid.Logging(log, reqDump, respDump)), mid.Recovery)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Infof("start echo server at %s", viper.GetString("app.host"))
			go func() {
				err := e.Start(viper.GetString("app.host"))
				if err != nil {
					log.Errorf("failed to start server, err:%+v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("start shuting down echo server")
			return e.Shutdown(ctx)
		},
	})

	return e, nil
}

func setupSwagger(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
