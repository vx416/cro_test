package wallet

import (
	"cro_test/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestWallet(t *testing.T) {
	suite.Run(t, new(WalletSuite))
}

type WalletSuite struct {
	testutil.TestSuite
}

func (s *WalletSuite) SetupTest() {
	s.DeleteTables("transactions", "wallets", "users")
}
