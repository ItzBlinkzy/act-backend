package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func ChatBotGroup(g *echo.Group) {
	g.POST("/chat-bot", controller.ChatBot)
}
