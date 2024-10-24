package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
)

var stockRepo = &repository.StockRepository{}

func ListBoughtStocks(ctx echo.Context) error {
	userIdParam := ctx.Param("userId")
	log.Printf("Received request to list stocks for user ID: %s", userIdParam)

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		log.Printf("Invalid user ID %s: %v", userIdParam, err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user ID"})
	}

	stocks, err := stockRepo.ListBoughtStocks(uint(userId))
	if err != nil {
		log.Printf("Failed to list stocks for user ID %d: %v", userId, err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list stocks"})
	}

	log.Printf("Stocks listed successfully for user ID %d", userId)
	return ctx.JSON(http.StatusOK, stocks)
}

func BuyStock(ctx echo.Context) error {
	var stock model.StockBuy
	if err := ctx.Bind(&stock); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid stock data"})
	}

	// Check if buying quantity is valid
	if stock.BuyingQuantity <= 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid buying quantity"})
	}

	err := stockRepo.BuyStock(stock)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to buy stock"})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "stock purchased"})
}

func UpdateStock(ctx echo.Context) error {
	stockIdParam := ctx.Param("stockId")
	stockId, err := strconv.Atoi(stockIdParam)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid stock ID"})
	}

	var update model.StockUpdate
	if err := ctx.Bind(&update); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input data"})
	}

	err = stockRepo.UpdateStock(uint(stockId), update)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update stock"})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "stock updated"})
}

// List logs by userId
func LogsBoughtStocksByUser(ctx echo.Context) error {
	userIdParam := ctx.Param("userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user ID"})
	}

	logs, err := stockRepo.GetLogsByUserId(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
	}
	return ctx.JSON(http.StatusOK, logs)
}

// List logs by stockId
func LogsBoughtStocksByStock(ctx echo.Context) error {
	stockIdParam := ctx.Param("stockId")
	stockId, err := strconv.Atoi(stockIdParam)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid stock ID"})
	}

	logs, err := stockRepo.GetLogsByStockId(stockId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
	}
	return ctx.JSON(http.StatusOK, logs)
}
