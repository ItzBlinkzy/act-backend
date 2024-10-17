package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func ClientGroup(g *echo.Group) {
	g.POST("/create-client", controller.CreateClient)
	g.PUT("/update-client/:clientId", controller.UpdateClient)
	g.DELETE("/delete-client/:clientId", controller.DeleteClient)
	g.GET("/client/:clientId", controller.GetClient)
	g.GET("/list-clients/:managerId", controller.GetAllClients)
	g.GET("/stocks-client/:clientId", controller.GetStocksOfClient)
}
