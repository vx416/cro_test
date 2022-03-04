package http

import (
	"cro_test/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func New(svc domain.Servicer) Handler {
	return Handler{
		svc: svc,
	}
}

type Handler struct {
	svc domain.Servicer
}

// GetUserID godoc
// @Summary      Ping
// @Tags         ping
// @Success      200  {string} string "ok"
// @Router       /ping [get]
func (h Handler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}
