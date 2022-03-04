package http

import (
	"cro_test/internal/model"
	"cro_test/pkg/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	Handler
}

type CreateWalletRequest struct {
	Currency string `json:"currency"`
}

// CreateWallet godoc
// @Summary      create wallet by email
// @Tags         wallet
// @Accept       json
// @Produce      json
// @security Bearer
// @Success      200  {object} model.Wallet
// @Router       /wallets [post]
func (h WalletHandler) CreateWallet(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}
	req := &CreateWalletRequest{}
	err = c.Bind(req)
	if err != nil {
		return err
	}

	wallet, err := h.svc.CreateWallet(ctx, &model.User{
		ID: userID,
	}, model.CurrencyValue(req.Currency))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

// CreateWallet godoc
// @Summary      get wallet by serial number
// @Tags         wallet
// @security Bearer
// @Param        serial   path      string  true  "Wallet Serail Number"
// @Success      200  {object} model.Wallet
// @Router       /wallets/{serial} [get]
func (h WalletHandler) GetWallet(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}
	serialNumber := c.Param("serial")
	wallet, err := h.svc.GetWalletBySerial(ctx, userID, serialNumber)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

// CreateWallet godoc
// @Summary      list wallets by userID
// @Tags         wallet
// @security Bearer
// @Success      200  {object} model.Wallet
// @Router       /wallets [get]
func (h WalletHandler) ListWallets(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}
	wallets, err := h.svc.ListWallets(ctx, model.ListWalletsOpts{
		UserID: userID,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallets)
}
