// config/routes.go
package config

import (
	v1 "github.com/itzblinkzy/act-backend/api/v1"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api")

	// Versioned API group
	v1Group := apiGroup.Group("/v1")
	v1.RegisterGroup(v1Group)
	v1.LoginGroup(v1Group)
	v1.LogoutGroup(v1Group)
	v1.UserGroup(v1Group)
	v1.StockGroup(v1Group)
	v1.ReviewGroup(v1Group)
	v1.ClientGroup(v1Group)
	v1.ChatBotGroup(v1Group)
	v1.PaymentGroup(v1Group)
}
