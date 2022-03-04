package middleware

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestToken_SignAndVeify(t *testing.T) {
	tcs := []struct {
		name   string
		exp    time.Time
		userID uint64
		secret []byte
		valid  bool
	}{
		{
			"create valid token", time.Now().Add(10 * time.Minute), 1, []byte("test"), true,
		},
		{
			"create invalid token", time.Now().Add(-10 * time.Minute), 2, []byte("test"), true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			claims := &Claims{
				UserID:    tc.userID,
				ExpiredAt: time.Now().Add(30 * time.Minute),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString(tc.secret)
			assert.NoError(t, err)
			assert.NotEmpty(t, signedToken)

			parsedClaims := &Claims{}

			tk, err := jwt.ParseWithClaims(signedToken, parsedClaims, func(token *jwt.Token) (interface{}, error) {
				return tc.secret, nil
			})
			assert.NoError(t, err)
			assert.True(t, tk.Valid)
			if !tc.valid {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tc.userID, parsedClaims.UserID)
			}
		})
	}
}
