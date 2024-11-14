package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func PaymentGroup(g *echo.Group) {
	g.POST("/add-credit/:userId", controller.AddCredit)
}
