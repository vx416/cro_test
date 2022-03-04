package middleware

import (
	"context"
	"cro_test/pkg/httperr"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const AuthenticationKey = "Authentication"

type (
	Claims struct {
		UserID    uint64    `json:"id"`
		ExpiredAt time.Time `json:"expiredAt"`
	}
	JWTTokenKey struct{}
)

func (c Claims) Valid() error {
	if c.ExpiredAt.Before(time.Now()) {
		return errors.New("expired")
	}
	return nil
}

func AuthJwtToken(secret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := extractToken(c.Request().Header, "Bearer", AuthenticationKey)
			if tokenStr == "" {
				return httperr.NewErr(http.StatusUnauthorized, "token not found")
			}
			claims, err := ParseWithClaims(tokenStr, secret)
			if err != nil {
				return err
			}
			ctx := context.WithValue(c.Request().Context(), JWTTokenKey{}, claims)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			return next(c)
		}
	}
}

func ParseWithClaims(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil || !token.Valid {
		return nil, httperr.WrapfErr(http.StatusUnauthorized, "token invalid", "err:%+v", err)
	}
	return claims, nil
}

func GetUserID(ctx context.Context) (uint64, error) {
	c, ok := ctx.Value(JWTTokenKey{}).(*Claims)
	if !ok {
		return 0, httperr.NewErr(http.StatusUnauthorized, "token not found")
	}
	return c.UserID, nil
}

func extractToken(header http.Header, schema, headerKey string) string {
	tokenHeader := header.Get(headerKey)
	if tokenHeader == "" {
		return ""
	}
	if len(schema)+1 > len(tokenHeader) || tokenHeader[:len(schema)] != schema {
		return ""
	}

	return tokenHeader[len(schema)+1:]
}
