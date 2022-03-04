package testutil

import (
	"context"
	"cro_test/internal/app"
	"cro_test/internal/delivery/http"
	"cro_test/internal/domain"
	"cro_test/internal/repository/mysql"
	"cro_test/pkg/config"
	"cro_test/pkg/container"
	"cro_test/pkg/logger"
	"cro_test/pkg/middleware"
	"database/sql"
	"time"

	"path/filepath"
	"runtime"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	gofactory "github.com/vx416/gogo-factory"
	"go.uber.org/fx"

	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	HttpHelper
	builder     *container.Builder
	Ctx         context.Context
	DB          *sql.DB
	MySQLDao    mysql.MySQLDao
	Service     domain.Servicer
	HttpHandler http.Handler
	App         *fx.App
	Echo        *echo.Echo
}

func (suite *TestSuite) SetupSuite() {
	suite.Ctx = context.Background()
	err := suite.SetupInstances()
	suite.Require().NoError(err, "init containers failed")
	err = suite.RunMigrations()
	suite.Require().NoError(err, "run migrations failed")
	suite.App, err = app.NewServerApp("test",
		fx.Populate(&suite.MySQLDao),
		fx.Populate(&suite.Service),
		fx.Populate(&suite.HttpHandler),
		fx.Populate(&suite.Echo),
		fx.Invoke(func(l logger.Logger) {
			suite.Ctx = l.Attach(suite.Ctx)
		}),
		fx.Provide(config.StartEchoServer),
		fx.Invoke(http.Routes),
	)
	suite.Require().NoError(err, "init app failed")
	err = suite.App.Start(suite.Ctx)
	suite.Require().NoError(err, "start app failed")
}

func (suite *TestSuite) SetupInstances() error {
	container, err := container.NewContainer()
	if err != nil {
		return err
	}
	suite.builder = container
	dbCfg, err := suite.builder.RunMysql("test-mysql", "cro_test", "13306")
	if err != nil {
		return err
	}

	db, err := dbCfg.GetDB()
	if err != nil {
		return err
	}
	suite.DB = db
	return nil
}

func (suite *TestSuite) RunMigrations() error {
	gofactory.Opt().SetDB(suite.DB, "mysql").SetTagProcess(gofactory.DBTagProcess)
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	migPath := filepath.Join(dir, "../../build/migrations")
	goose.SetDialect("mysql")
	return goose.Run("up", suite.DB, migPath)
}

func (suite *TestSuite) DeleteTables(tables ...string) {
	var err error
	for _, table := range tables {
		_, dbErr := suite.DB.Exec("DELETE FROM " + table)
		if dbErr != nil {
			err = dbErr
		}
	}
	suite.Require().NoError(err, "delete tables failed")
}

func (suite *TestSuite) GetFakeToken(userID uint64) string {
	claims := middleware.Claims{
		UserID:    userID,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(config.GetJwtSecret())
	suite.Require().NoError(err, "get fake token failed")
	return token
}
