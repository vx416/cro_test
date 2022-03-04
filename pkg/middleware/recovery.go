package middleware

import (
	"cro_test/pkg/logger"
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
)

func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				trace := make([]byte, 4096)
				runtime.Stack(trace, true)
				err, ok := r.(error)
				if !ok {
					if err == nil {
						err = fmt.Errorf("%v", r)
					} else {
						err = fmt.Errorf("%v", err)
					}
				}
				logger.Ctx(c.Request().Context()).Err(err).Field("stack_err", string(trace)).Error("panice!!!!!!!")
				_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": http.StatusText(http.StatusInternalServerError)})
			}
		}()
		return next(c)
	}
}
