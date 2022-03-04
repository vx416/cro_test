package transaction

import (
	"cro_test/internal/model"
	"cro_test/internal/test/factory"
	"sync"

	"github.com/shopspring/decimal"
)

func (s *TxSuite) TestSvc_Transfer() {
	u1 := factory.User.HasWallets(factory.Wallet.Currency("twd").Amount("100"), 1).MustInsert().(*model.User)
	u2 := factory.User.HasWallets(factory.Wallet.Currency("twd").Amount("0"), 1).MustInsert().(*model.User)
	tcs := []struct {
		name       string
		userID     uint64
		fromSerial string
		toSerial   string
		delta      decimal.Decimal
		hasErr     bool
		fromAmount string
		toAmount   string
	}{
		{
			"transfer_ok",
			u1.ID, u1.Wallets[0].SerialNumber, u2.Wallets[0].SerialNumber, decimal.NewFromInt(50),
			false, "50", "50",
		},
		{
			"transfer_insufficient",
			u1.ID, u1.Wallets[0].SerialNumber, u2.Wallets[0].SerialNumber, decimal.NewFromInt(100),
			true, "50", "50",
		},
		{
			"transfer_user_invalid",
			u2.ID, u1.Wallets[0].SerialNumber, u2.Wallets[0].SerialNumber, decimal.NewFromInt(100),
			true, "50", "50",
		},
		{
			"transfer_to_serial_number_not_found",
			u2.ID, u1.Wallets[0].SerialNumber, "test_not_found", decimal.NewFromInt(100),
			true, "50", "0",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			tx, err := s.Service.Transfer(s.Ctx, tc.userID, tc.fromSerial, tc.toSerial, tc.delta)

			if tc.hasErr {
				s.Error(err)
				fromWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.fromSerial)
				s.Require().NoError(err, "get wallet failed")
				s.Equal(fromWallet.Amount.String(), tc.fromAmount)
				toWallet, _ := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.toSerial)
				s.Equal(toWallet.Amount.String(), tc.toAmount)
				return
			}
			s.Require().NoError(err, "update wallet failed")
			fromWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.fromSerial)
			s.Require().NoError(err, "get wallet failed")
			s.Equal(fromWallet.Amount.String(), tc.fromAmount)
			s.Equal(tx.FromWalletBalance.String(), tc.fromAmount)
			toWallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.toSerial)
			s.Require().NoError(err, "get wallet failed")
			s.Equal(tx.ToWalletBalance.String(), tc.fromAmount)
			s.Equal(toWallet.Amount.String(), tc.toAmount)
		})
	}
}

func (s *TxSuite) TestSvc_Withdraw() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("150"), 1).MustInsert().(*model.User)
	u2 := factory.User.HasWallets(factory.Wallet.Amount("0"), 1).MustInsert().(*model.User)
	tcs := []struct {
		name         string
		userID       uint64
		serial       string
		delta        decimal.Decimal
		hasErr       bool
		expectAmount string
	}{
		{
			"withdraw_ok",
			u1.ID, u1.Wallets[0].SerialNumber, decimal.NewFromInt(100), false, "50",
		},
		{
			"withdraw_insufficient",
			u1.ID, u1.Wallets[0].SerialNumber, decimal.NewFromInt(100), true, "0",
		},
		{
			"withdraw_invalid_user_id",
			u2.ID, u1.Wallets[0].SerialNumber, decimal.NewFromInt(100), true, "0",
		},
		{
			"withdraw_insufficient",
			u2.ID, u2.Wallets[0].SerialNumber, decimal.NewFromInt(100), true, "0",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			tx, err := s.Service.Withdraw(s.Ctx, tc.userID, tc.serial, tc.delta)
			if tc.hasErr {
				s.Error(err)
				return
			}
			wallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.serial)
			s.Require().NoError(err, "get wallet failed")
			s.Equal(wallet.Amount.String(), tc.expectAmount)
			s.Require().NoError(err, "update wallet failed")
			s.Equal(tx.FromWalletBalance.String(), tc.expectAmount)
		})
	}
}

func (s *TxSuite) TestSvc_Deposit() {
	u1 := factory.User.HasWallets(factory.Wallet.Amount("0"), 1).MustInsert().(*model.User)
	u2 := factory.User.HasWallets(factory.Wallet.Amount("0"), 1).MustInsert().(*model.User)
	tcs := []struct {
		name         string
		userID       uint64
		serial       string
		delta        decimal.Decimal
		hasErr       bool
		expectAmount string
	}{
		{
			"deposit_ok",
			u1.ID, u1.Wallets[0].SerialNumber, decimal.NewFromInt(100), false, "100",
		},
		{
			"deposit_invalid_user_id",
			u2.ID, u1.Wallets[0].SerialNumber, decimal.NewFromInt(100), true, "0",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			tx, err := s.Service.Deposit(s.Ctx, tc.userID, tc.serial, tc.delta)
			if tc.hasErr {
				s.Error(err)
				return
			}
			wallet, err := s.MySQLDao.GetWalletBySerial(s.Ctx, tc.serial)
			s.Require().NoError(err, "get wallet failed")
			s.Equal(wallet.Amount.String(), tc.expectAmount)
			s.Require().NoError(err, "update wallet failed")
			s.Equal(tx.ToWalletBalance.String(), tc.expectAmount)
		})
	}
}

func (s *TxSuite) TestSvc_Concurrency() {
	u1 := factory.User.HasWallets(factory.Wallet.Currency("twd").Amount("2000"), 1).MustInsert().(*model.User)
	u2 := factory.User.HasWallets(factory.Wallet.Currency("twd").Amount("0"), 1).MustInsert().(*model.User)
	w1 := u1.Wallets[0]
	w2 := u2.Wallets[0]
	wg := sync.WaitGroup{}

	transfer := decimal.NewFromInt(1)
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Service.Transfer(s.Ctx, u1.ID, w1.SerialNumber, w2.SerialNumber, transfer)
		}()
	}

	withdraw := decimal.NewFromInt(2)
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Service.Withdraw(s.Ctx, u1.ID, w1.SerialNumber, withdraw)
		}()
	}

	wg.Wait()

	ww1, err := s.MySQLDao.GetWalletBySerial(s.Ctx, w1.SerialNumber)
	s.Require().NoError(err)
	ww2, err := s.MySQLDao.GetWalletBySerial(s.Ctx, w2.SerialNumber)
	s.Require().NoError(err)
	s.Equal("1100", ww1.Amount.String())
	s.Equal("300", ww2.Amount.String())
}
