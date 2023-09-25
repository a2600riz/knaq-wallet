package groups

import (
	"github.com/labstack/echo/v4"
	"knaq-wallet/tools/security"
)

const (
	ServiceName = "wallet"
)

var (
	securePath   *echo.Group
	internalPath *echo.Group
	globalPath   *echo.Group
	adminPath    *echo.Group
)

func NewGroup(server *echo.Echo) {
	globalPath = server.Group("/" + ServiceName)
	configureSecurePath(globalPath)
	configureInternalPath(server.Group(""))

	registerHandler()
}
func configureSecurePath(group *echo.Group) {
	securePath = group.Group("/api/v1")
	securePath.Use(security.ProcessUserClaim)
}
func configureInternalPath(group *echo.Group) {
	internalPath = group.Group("/internal/" + ServiceName + "/api/v1")
}
