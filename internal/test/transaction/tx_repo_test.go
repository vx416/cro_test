package transaction

import (
	"cro_test/internal/model"
	"cro_test/internal/repository/mysql"
	"cro_test/internal/test/factory"

	"github.com/shopspring/decimal"
)

func (s *TxSuite) TestRepo_UpdateWalletDelta() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("0"), 1).MustInsert().(*model.User)
	tcs := []struct {
		serial       string
		delta        decimal.Decimal
		hasErr       bool
		expectAmount string
	}{
		{
			u1.Wallets[0].SerialNumber, decimal.NewFromInt(100), false, "100",
		},
		{
			u1.Wallets[0].SerialNumber, decimal.NewFromInt(-100), false, "0",
		},
		{
			u1.Wallets[0].SerialNumber, decimal.NewFromInt(-100), true, "0",
		},
	}

	for _, tc := range tcs {
		s.Run("update_delta", func() {
			wallet, err := s.MySQLDao.UpdateWalletAmountDelta(s.Ctx, tc.serial, tc.delta)
			if tc.hasErr {
				s.Error(err)
				s.ErrorIs(err, mysql.ErrNotAffected)
				return
			}
			s.Require().NoError(err, "update wallet failed")
			s.Equal(wallet.SerialNumber, tc.serial)
			s.Equal(wallet.Amount.String(), tc.expectAmount)
		})
	}
}
