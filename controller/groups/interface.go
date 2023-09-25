package groups

import "github.com/labstack/echo/v4"

type Handler interface {
	Register(group *echo.Group)
}
