package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func StockGroup(g *echo.Group) {
	g.GET("/list-bought-stocks/:userId", controller.ListBoughtStocks)
	g.POST("/buy-stock", controller.BuyStock)
	g.PUT("/update-stock/:stockId", controller.UpdateStock)
	g.GET("/list-logs-bought-stock-user/:userId", controller.LogsBoughtStocksByUser)
	g.GET("/list-logs-bought-stock/:stockId", controller.LogsBoughtStocksByStock)
}
