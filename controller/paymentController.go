package controller

import (
	"net/http"
	"strconv"

	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
)

func AddCredit(c echo.Context) error {
	userID := c.Param("userId")
	creditStr := c.QueryParam("credit")

	credit, err := strconv.ParseFloat(creditStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid credit amount"})
	}

	if err := repository.PaymentRepo.AddCreditToUser(userID, credit); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Credit added successfully"})
}
