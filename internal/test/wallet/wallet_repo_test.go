package wallet

import (
	"cro_test/internal/model"
	"cro_test/internal/test/factory"
	"database/sql"
)

func (s *WalletSuite) TestRepo_GetWalletBySerial() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("0"), 1).MustInsert().(*model.User)
	tcs := []struct {
		serial string
		hasErr bool
	}{
		{
			"not_found", true,
		},
		{
			u1.Wallets[0].SerialNumber, false,
		},
	}

	for _, tc := range tcs {
		s.Run("getwallet", func() {
			wallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.serial)
			if tc.hasErr {
				s.Error(err)
				s.ErrorIs(err, sql.ErrNoRows)
				return
			}
			s.Require().NoError(err, "get wallet failed")
			s.Equal(wallet.SerialNumber, tc.serial)
		})
	}
}

func (s *WalletSuite) TestRepo_ListWallets() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("0"), 3).MustInsert().(*model.User)
	u2 := factory.User.HasWallets(factory.Wallet.Amount("0"), 2).MustInsert().(*model.User)
	tcs := []struct {
		opts      model.ListWalletsOpts
		hasErr    bool
		expectCnt int
	}{
		{
			model.ListWalletsOpts{
				UserID: u1.ID,
			}, false, 3,
		},
		{
			model.ListWalletsOpts{
				SerilaNumbers: []string{u1.Wallets[0].SerialNumber, u2.Wallets[0].SerialNumber},
			}, false, 2,
		},
		{
			model.ListWalletsOpts{}, false, 5,
		},
		{
			model.ListWalletsOpts{
				UserID: 123,
			}, false, 0,
		},
	}

	for _, tc := range tcs {
		s.Run("getwallet", func() {
			wallets, err := s.MySQLDao.ListWallets(s.Ctx, tc.opts)
			if tc.hasErr {
				s.Error(err)

				return
			}
			s.Require().NoError(err, "get wallet failed")
			s.Len(wallets, tc.expectCnt)
			if tc.opts.UserID > 0 {
				for _, wallet := range wallets {
					s.Equal(wallet.UserID, tc.opts.UserID)
				}
			}
			if len(tc.opts.SerilaNumbers) > 0 {
				s.Len(wallets, len(tc.opts.SerilaNumbers))
			}
		})
	}
}
