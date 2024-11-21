package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
)

func CreateClient(c echo.Context) error {
	var req struct {
		ManagerID   uint   `json:"manager_id"`
		CompanyName string `json:"company_name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	client := model.Client{
		CompanyName: req.CompanyName,
	}
	clientID, err := repository.CreateClient(client)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create client"})
	}

	association := model.ClientManagerAssociation{
		ManagerId: req.ManagerID,
		ClientId:  clientID,
	}
	err = repository.CreateAssociation(association)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create association: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "client created successfully"})
}

func UpdateClient(c echo.Context) error {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid client ID"})
	}

	var req struct {
		CompanyName string `json:"company_name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err = repository.UpdateClient(uint(clientId), req.CompanyName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update client"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "client updated successfully"})
}

func DeleteClient(c echo.Context) error {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid client ID"})
	}

	if err := repository.SoftDeleteClient(uint(clientId)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete client"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "client deleted successfully"})
}

func GetClient(c echo.Context) error {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid client ID"})
	}

	client, err := repository.GetClient(uint(clientId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "client not found"})
	}

	return c.JSON(http.StatusOK, client)
}

func GetAllClients(c echo.Context) error {
	managerIdParam := c.Param("managerId")
	managerId, err := strconv.ParseUint(managerIdParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid manager ID"})
	}

	clients, err := repository.GetAllClients(uint(managerId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve clients"})
	}

	var response []map[string]interface{}

	for _, client := range clients {
		stocks, err := repository.GetStocksOfClient(client.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve stocks"})
		}

		response = append(response, map[string]interface{}{
			"client": client,
			"stocks": stocks,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func GetStocksOfClient(c echo.Context) error {
	clientIdParam := c.Param("clientId")
	clientId, err := strconv.ParseUint(clientIdParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid client ID"})
	}

	stocks, err := repository.GetStocksOfClient(uint(clientId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve stocks"})
	}

	return c.JSON(http.StatusOK, stocks)
}
