package auth

import (
	"cro_test/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

type AuthSuite struct {
	testutil.TestSuite
}

func (s *AuthSuite) SetupTest() {
	s.DeleteTables("transactions", "wallets", "users")
}
