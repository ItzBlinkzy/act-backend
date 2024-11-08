package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/itzblinkzy/act-backend/model"
	"github.com/labstack/echo/v4"
)

func ChatBot(c echo.Context) error {
	var requestData model.RequestDataChatBot
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	prompt := fmt.Sprintf("Here is the conversation history: %s\n\nUser: %s\n\nAI:", requestData.Context, requestData.Message)
	requestBody := map[string]interface{}{
		"model": "ft:gpt-4o-mini-2024-07-18:personal::ARKkt0BQ",
		"messages": []map[string]string{
			{"role": "system", "content": "This assistant is designed to help users with inquiries about cryptocurrencies, stock market trends, and any questions related to the company's website or mobile application. The assistant provides guidance on navigating the platform, understanding market data, and making informed decisions based on current trends. Please note: all advice is for informational purposes only and does not constitute financial or investment advice."},
			{"role": "user", "content": prompt},
		},
	}

	client := resty.New()

	apiKey := os.Getenv("CHAT_BOT_API_KEY")
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to contact AI service")
	}

	var chatGPTResp model.ChatGPTResponse
	if err := json.Unmarshal(resp.Body(), &chatGPTResp); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse AI response")
	}

	return c.JSON(http.StatusOK, chatGPTResp)
}
