package controller

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"knaq-wallet/config"
	"knaq-wallet/controller/erra"
	"knaq-wallet/controller/groups"
	"knaq-wallet/tools/valid"
	"os"
	"strings"
)

var (
	server *echo.Echo
)

func Start() {

	server = echo.New()
	server.Use(middleware.Recover())
	server.Use(middleware.Secure())
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(ctx echo.Context) bool {
			if strings.Contains(ctx.Request().RequestURI, "/health") {
				return true
			}
			return false
		},
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           os.Stdout,
	}))
	server.HTTPErrorHandler = erra.HTTPErrorHandler
	server.Validator = &valid.CustomValidator{Validator: validator.New()}

	groups.NewGroup(server)

	switch config.Config.GetStage() {
	case config.PROD:
		server.Logger.SetLevel(log.ERROR)
	default:
		server.Logger.SetLevel(log.DEBUG)
	}

	server.Logger.Fatal(server.Start(fmt.Sprintf(":%s", config.Config.GetPort())))
}
