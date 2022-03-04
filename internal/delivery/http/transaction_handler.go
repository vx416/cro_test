package http

import (
	"cro_test/internal/model"
	"cro_test/pkg/httperr"
	"cro_test/pkg/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	FromWalletSerial string `json:"fromWalletSerial"`
	ToWalletSerial   string `json:"toWalletSerial"`
	Amount           string `json:"amount"`
}

func NewTxResp(tx model.Transaction) TransactionResponse {
	switch tx.Kind {
	case model.TxKindTransfer, model.TxKindWithdraw:
		return TransactionResponse{
			FromWalletBalance: tx.FromWalletBalance.String(),
		}
	case model.TxKindDeposit:
		return TransactionResponse{
			ToWalletBalance: tx.ToWalletBalance.String(),
		}
	}

	return TransactionResponse{}
}

type TransactionResponse struct {
	FromWalletBalance string `json:"fromWalletBalance,omitempty"`
	ToWalletBalance   string `json:"toWalletBalance,omitempty"`
}

func (req TransactionRequest) GetAmount() (decimal.Decimal, error) {
	d, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return decimal.Decimal{}, httperr.WrapfErr(http.StatusUnprocessableEntity, "amount is invalid", "new decimal from %s failed", req.Amount)
	}
	return d, nil
}

type TransactionHandler struct {
	Handler
}

// CreateTransfer godoc
// @Summary      create transfer
// @Tags         transaction
// @Accept       json
// @Produce      json
// @security Bearer
// @Param        TranscationRequest   body      TransactionRequest  true  "transaction request"
// @Success      200  {object} TransactionResponse
// @failure      422  {string} string "wallet balance insufficient"
// @failure      401  {string} string "not wallet owner"
// @Router       /transfer [post]
func (h TransactionHandler) CreateTransfer(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}
	req := &TransactionRequest{}
	err = c.Bind(req)
	if err != nil {
		return err
	}
	a, err := req.GetAmount()
	if err != nil {
		return err
	}

	tx, err := h.svc.Transfer(ctx, userID, req.FromWalletSerial, req.ToWalletSerial, a)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewTxResp(tx))
}

// CreateDeposit godoc
// @Summary      create deposit
// @Tags         transaction
// @Accept       json
// @Produce      json
// @security Bearer
// @Param        TranscationRequest   body      TransactionRequest  true  "only toWalletSerial and amount"
// @Success      200  {object} TransactionResponse
// @failure      422  {string} string "wallet balance insufficient"
// @failure      401  {string} string "not wallet owner"
// @Router       /deposit [post]
func (h TransactionHandler) CreateDeposit(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}
	req := &TransactionRequest{}
	err = c.Bind(req)
	if err != nil {
		return err
	}
	a, err := req.GetAmount()
	if err != nil {
		return err
	}
	tx, err := h.svc.Deposit(ctx, userID, req.ToWalletSerial, a)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewTxResp(tx))
}

// CreateWithdraw godoc
// @Summary      create deposit
// @Tags         transaction
// @Accept       json
// @Produce      json
// @security Bearer
// @Param        TranscationRequest   body      TransactionRequest  true  "only fromWalletSerial and amount"
// @Success      200  {object} TransactionResponse
// @failure      422  {string} string "wallet balance insufficient"
// @failure      401  {string} string "not wallet owner"
// @Router       /withdraw [post]
func (h TransactionHandler) CreateWithdraw(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}
	req := &TransactionRequest{}
	err = c.Bind(req)
	if err != nil {
		return err
	}
	a, err := req.GetAmount()
	if err != nil {
		return err
	}
	tx, err := h.svc.Withdraw(ctx, userID, req.FromWalletSerial, a)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewTxResp(tx))
}
