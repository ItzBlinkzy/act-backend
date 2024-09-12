package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func UserGroup(g *echo.Group) {
	g.GET("/users", controller.GetAllUsers)
	g.GET("/me", controller.GetCurrentUser)
	g.PUT("/update-user/:id", controller.UpdateUser)
}
