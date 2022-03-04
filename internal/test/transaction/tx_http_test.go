package transaction

import (
	"cro_test/internal/model"
	"cro_test/internal/test/factory"
)

func (s *TxSuite) TestHttp_Transfer() {
	u1 := factory.User.HasWallets(factory.Wallet.Currency("twd").Amount("50"), 1).MustInsert().(*model.User)
	u2 := factory.User.HasWallets(factory.Wallet.Currency("twd").Amount("50"), 1).MustInsert().(*model.User)
	req, resp := s.HttpHelper.BuildRequest("POST", "/api/v1/transfer", map[string]interface{}{
		"fromWalletSerial": u1.Wallets[0].SerialNumber,
		"toWalletSerial":   u2.Wallets[0].SerialNumber,
		"amount":           "10",
	})
	s.HttpHelper.SetBearToken(req, s.GetFakeToken(u1.ID))
	s.Echo.Server.Handler.ServeHTTP(resp, req)
	fromWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, u1.Wallets[0].SerialNumber)
	s.Require().NoError(err)
	s.Equal("40", fromWallet.Amount.String())
	toWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, u2.Wallets[0].SerialNumber)
	s.Require().NoError(err)
	s.Equal("60", toWallet.Amount.String())
}

func (s *TxSuite) TestHttp_Withdraw() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("50"), 1).MustInsert().(*model.User)
	req, resp := s.HttpHelper.BuildRequest("POST", "/api/v1/withdraw", map[string]interface{}{
		"fromWalletSerial": u1.Wallets[0].SerialNumber,
		"amount":           "10",
	})
	s.HttpHelper.SetBearToken(req, s.GetFakeToken(u1.ID))
	s.Echo.Server.Handler.ServeHTTP(resp, req)
	fromWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, u1.Wallets[0].SerialNumber)
	s.Require().NoError(err)
	s.Equal("40", fromWallet.Amount.String())
}

func (s *TxSuite) TestHttp_Deposit() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("50"), 1).MustInsert().(*model.User)
	req, resp := s.HttpHelper.BuildRequest("POST", "/api/v1/deposit", map[string]interface{}{
		"toWalletSerial": u1.Wallets[0].SerialNumber,
		"amount":         "10",
	})
	s.HttpHelper.SetBearToken(req, s.GetFakeToken(u1.ID))
	s.Echo.Server.Handler.ServeHTTP(resp, req)
	fromWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, u1.Wallets[0].SerialNumber)
	s.Require().NoError(err)
	s.Equal("60", fromWallet.Amount.String())
}
