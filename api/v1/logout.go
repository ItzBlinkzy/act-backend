package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func LogoutGroup(g *echo.Group) {
	g.POST("/logout", controller.Logout)
}
