package http

import (
	"cro_test/pkg/config"
	"cro_test/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo, handler Handler) {
	{
		v1Api := e.Group("/api/v1")
		v1Api.GET("/ping", handler.Ping)
		walletRoutes(v1Api, WalletHandler{handler})
		transactionRoutes(v1Api, TransactionHandler{handler})
		authRoutes(v1Api, AuthHandler{handler})
	}
}

func walletRoutes(g *echo.Group, handler WalletHandler) {
	g.POST("/wallets", handler.CreateWallet, middleware.AuthJwtToken(config.GetJwtSecret()))
	g.GET("/wallets/:serial", handler.GetWallet, middleware.AuthJwtToken(config.GetJwtSecret()))
	g.GET("/wallets", handler.ListWallets, middleware.AuthJwtToken(config.GetJwtSecret()))
}

func transactionRoutes(g *echo.Group, handler TransactionHandler) {
	g.POST("/transfer", handler.CreateTransfer, middleware.AuthJwtToken(config.GetJwtSecret()))
	g.POST("/deposit", handler.CreateDeposit, middleware.AuthJwtToken(config.GetJwtSecret()))
	g.POST("/withdraw", handler.CreateWithdraw, middleware.AuthJwtToken(config.GetJwtSecret()))
}

func authRoutes(g *echo.Group, handler AuthHandler) {
	g.POST("/auth", handler.Login)
	g.POST("/signup", handler.Signup)
}
