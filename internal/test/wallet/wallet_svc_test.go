package wallet

import (
	"cro_test/internal/model"
	"cro_test/internal/test/factory"
)

func (s *WalletSuite) TestSvc_CreateWallet() {
	users := factory.User.MustInsertN(2).([]*model.User)
	tcs := []struct {
		name     string
		user     *model.User
		hasErr   bool
		currency model.CurrencyValue
	}{
		{
			"create_wallet_ok",
			users[0], false,
			"twd",
		},
		{
			"create_wallet_ok",
			users[1], false,
			"usd",
		},
		{
			"create_wallet_currency_duplicated",
			users[1], true,
			"usd",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			wallet, err := s.Service.CreateWallet(s.Ctx, tc.user, tc.currency)
			if tc.hasErr {
				s.Error(err)
				return
			}
			s.Require().NoError(err, "create wallet failed")
			s.Equal(wallet.Amount.String(), "0")
			s.Equal(wallet.UserID, tc.user.ID)
		})
	}
}
