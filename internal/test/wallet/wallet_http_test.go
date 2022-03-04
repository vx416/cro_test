package wallet

import (
	"cro_test/internal/model"
	"cro_test/internal/test/factory"
)

func (s *WalletSuite) TestHttp_CreateWallet() {
	user := factory.User.MustInsert().(*model.User)

	req, resp := s.HttpHelper.BuildRequest("POST", "/api/v1/wallets", map[string]interface{}{
		"currency": "twd",
	})
	s.HttpHelper.SetBearToken(req, s.GetFakeToken(user.ID))
	s.Echo.Server.Handler.ServeHTTP(resp, req)
	wallet := model.Wallet{}
	err := s.HttpHelper.GetResponseData(resp, &wallet)
	s.Require().NoError(err)
	s.NotEmpty(wallet.SerialNumber)
	s.NotZero(wallet.UserID)
	s.Equal("0", wallet.Amount.String())
}

func (s *WalletSuite) TestHttp_GetWallet() {
	user := factory.User.HasWallets(factory.Wallet.Amount("100"), 1).MustInsert().(*model.User)
	req, resp := s.HttpHelper.BuildRequest("GET", "/api/v1/wallets/"+user.Wallets[0].SerialNumber, nil)
	s.HttpHelper.SetBearToken(req, s.GetFakeToken(user.ID))
	s.Echo.Server.Handler.ServeHTTP(resp, req)
	wallet := model.Wallet{}
	err := s.HttpHelper.GetResponseData(resp, &wallet)
	s.Require().NoError(err)
	s.Equal(user.Wallets[0].SerialNumber, wallet.SerialNumber)
	s.Equal(user.ID, wallet.UserID)
	s.Equal("100", wallet.Amount.String())
}

func (s *WalletSuite) TestHttp_GetWallets() {
	user := factory.User.HasWallets(factory.Wallet.Amount("100"), 2).MustInsert().(*model.User)
	req, resp := s.HttpHelper.BuildRequest("GET", "/api/v1/wallets", nil)
	s.HttpHelper.SetBearToken(req, s.GetFakeToken(user.ID))

	s.Echo.Server.Handler.ServeHTTP(resp, req)
	wallets := []model.Wallet{}
	err := s.HttpHelper.GetResponseData(resp, &wallets)
	s.Require().NoError(err)

	s.Len(wallets, 2)
	for i, wallet := range wallets {
		s.Equal(user.Wallets[i].SerialNumber, wallet.SerialNumber)
		s.Equal(user.ID, wallet.UserID)
		s.Equal("100", wallet.Amount.String())
	}

}
