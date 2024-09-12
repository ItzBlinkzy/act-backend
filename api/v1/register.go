package v1

import (
	"github.com/itzblinkzy/act-backend/controller"

	"github.com/labstack/echo/v4"
)

func RegisterGroup(g *echo.Group) {
	g.POST("/register-company-first-user", controller.Register)
}
