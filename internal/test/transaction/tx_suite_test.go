package transaction

import (
	"cro_test/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTransaction(t *testing.T) {
	suite.Run(t, new(TxSuite))
}

type TxSuite struct {
	testutil.TestSuite
}

func (s *TxSuite) SetupTest() {
	s.DeleteTables("transactions", "wallets", "users")
}
