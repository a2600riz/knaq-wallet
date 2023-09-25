package erra

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}
	var msg string

	if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
		code = http.StatusBadRequest
	}

	switch code {
	case http.StatusBadRequest, http.StatusUnauthorized:
		if he != nil {
			msgMap, ok := he.Message.(map[string]interface{})
			if ok {
				msgErr, ok := msgMap["error"].(error)
				if ok {
					msg = msgErr.Error()
					break
				}
			}
		}
		msg = gorm.ErrRecordNotFound.Error()
	case http.StatusInternalServerError:
		msg = "internal_server_error"
	}

	err = c.JSON(code, map[string]interface{}{
		"@timestamp": time.Now(),
		"message":    msg,
		"status":     code,
		"uri":        c.Request().RequestURI,
	})
	if err != nil {
		return
	}
}

func Error(err error) error {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Errorf(fmt.Sprintf("[ERROR] %s [%s:%d]\n", err, fn, line))
}
func BadRequest(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
		"message": "client_bad_request",
		"error":   err,
	})
}
func PreconditionRequired(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusPreconditionRequired, message)
}
func Unauthorized(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
		"message": "unauthorized",
		"error":   err,
	})
}
