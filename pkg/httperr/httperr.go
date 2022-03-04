package httperr

import (
	"cro_test/pkg/logger"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func ParseHttpErr(err error) (HttpError, bool) {
	err = errors.Cause(err)
	httErr, ok := err.(HttpError)
	if ok {
		return httErr, true
	}
	return HttpError{
		StatusCode: http.StatusInternalServerError,
		ErrResp:    make(map[string]interface{}),
	}, false
}

func NewErr(status int, errMsg string) error {
	if errMsg == "" {
		errMsg = http.StatusText(status)
	}

	return HttpError{
		StatusCode: status,
		ErrResp: map[string]interface{}{
			"error": errMsg,
		},
	}
}

func WrapfErr(status int, errMsg string, wrapMsg string, args ...interface{}) error {
	return errors.Wrapf(NewErr(status, errMsg), wrapMsg, args...)
}

type HttpError struct {
	StatusCode int
	ErrResp    map[string]interface{}
}

func (er HttpError) Error() string {
	var errMsg interface{} = http.StatusText(er.StatusCode)
	if er.ErrResp["error"] != nil {
		errMsg = er.ErrResp["error"]
	}
	return fmt.Sprintf("status %d, msg:%+v", er.StatusCode, errMsg)
}

func (er HttpError) Resp(c echo.Context) {
	if er.StatusCode == 0 {
		er.StatusCode = http.StatusInternalServerError
	}
	if er.ErrResp["error"] == nil {
		er.ErrResp["error"] = http.StatusText(er.StatusCode)
	}

	c.JSON(er.StatusCode, er.ErrResp)
}

func EchoErrHandle(err error, c echo.Context) {
	l := logger.Ctx(c.Request().Context())
	httpErr, ok := ParseHttpErr(err)
	if !ok {
		l.Warn("error is not http error use internal error as response")
	}
	httpErr.Resp(c)
	if c.Response().Status >= 400 {
		l.Field("error", err.Error()).Warnf("occur error, err:%+v", err)
	} else if c.Response().Status >= 500 {
		l.Field("error", err.Error()).Errorf("occur error, err:%+v", err)
	}
}
