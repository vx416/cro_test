package http

import (
	"cro_test/pkg/config"
	"cro_test/pkg/middleware"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string    `json:"accessToken"`
	ExpiredAt   time.Time `json:"expiredAt"`
}

type AuthHandler struct {
	Handler
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AuthRequest   body  AuthRequest  true "auth"
// @Success      200  {model} TokenResponse
// @Router       /auth [post]
func (h AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	req := AuthRequest{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}
	claims := &middleware.Claims{
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(config.GetJwtSecret())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken: token,
		ExpiredAt:   claims.ExpiredAt,
	})
}

// Singup godoc
// @Summary      SingUp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AuthRequest   body  AuthRequest  true "auth"
// @Success      200  {string} string "ok"
// @Router       /signup [post]
func (h AuthHandler) Signup(c echo.Context) error {
	ctx := c.Request().Context()
	req := AuthRequest{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	_, err = h.svc.SignUp(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}
