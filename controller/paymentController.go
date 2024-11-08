package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
	paypal "github.com/plutov/paypal/v3"
)

func newPayPalClient() (*paypal.Client, error) {
	clientID := os.Getenv("PAYPAL_CLIENT_ID")
	secret := os.Getenv("PAYPAL_SECRET")
	client, err := paypal.NewClient(clientID, secret, paypal.APIBaseSandBox)
	if err != nil {
		return nil, err
	}

	// This method handles the OAuth2 token under the hood
	_, err = client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	client.SetLog(os.Stdout) // Set log to standard output
	return client, nil
}

func AddCredit(c echo.Context) error {
	client, err := newPayPalClient()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create PayPal client")
	}

	amount := c.FormValue("amount")
	userID := c.Param("userId")

	// Adjusted call to match expected function signature
	order, err := client.CreateOrder(paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{
		{
			ReferenceID: userID,
			Amount: &paypal.PurchaseUnitAmount{
				Currency: "USD",
				Value:    amount,
			},
			Description: "Add credit to account",
		},
	}, nil, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create PayPal order: "+err.Error())
	}

	// Redirect user to the approval URL
	for _, link := range order.Links {
		if link.Rel == "approve" {
			return c.Redirect(http.StatusFound, link.Href)
		}
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "No approval link provided by PayPal")
}

func CapturePayment(c echo.Context) error {
	client, err := newPayPalClient()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create PayPal client")
	}

	orderID := c.QueryParam("orderID")

	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No order ID provided")
	}

	capturedOrder, err := client.CaptureOrder(orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to capture PayPal order: "+err.Error())
	}

	// Accessing captured amount details typically from the transaction captures
	var totalAmount float64
	for _, capture := range capturedOrder.PurchaseUnits[0].Payments.Captures {
		amt, err := strconv.ParseFloat(capture.Amount.Value, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid amount format: "+err.Error())
		}
		totalAmount += amt
	}

	userID, err := strconv.Atoi(capturedOrder.PurchaseUnits[0].ReferenceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if err := repository.UserRepo.UpdateUserCredit(userID, totalAmount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user credit: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Payment captured and credit added successfully",
		"userID":  userID,
		"amount":  totalAmount,
	})
}

// Assuming the correct paths from the logged JSON
type PayPalWebhookEvent struct {
	ID        string `json:"id"`
	EventType string `json:"event_type"`
	Resource  struct {
		Capture struct {
			ID     string `json:"id"`
			Amount struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
			FinalCapture bool `json:"final_capture"`
		} `json:"capture"`
		CustomID string `json:"custom_id"` // This must match the actual JSON key for custom ID
	} `json:"resource"`
}

type PayPalResource struct {
	Amount struct {
		Value string `json:"value"`
	} `json:"amount"`
	CustomID string `json:"custom_id"` // Custom ID might be your reference ID
}

func PayPalWebhook(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body")
	}
	defer c.Request().Body.Close()

	var notification PayPalWebhookEvent
	err = json.Unmarshal(body, &notification)
	if err != nil {
		log.Printf("Error parsing webhook JSON: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse webhook JSON")
	}

	log.Printf("Webhook received: %s, Type: %s", notification.ID, notification.EventType)

	// Handle different event types
	switch notification.EventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		return handlePaymentCaptureCompleted(notification, c)
	default:
		log.Printf("Unhandled event type: %s", notification.EventType)
		return echo.NewHTTPError(http.StatusNotImplemented, "Event type not handled")
	}
}

func handlePaymentCaptureCompleted(notification PayPalWebhookEvent, c echo.Context) error {
	userID, err := strconv.Atoi(notification.Resource.CustomID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	amount, err := strconv.ParseFloat(notification.Resource.Capture.Amount.Value, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid amount format")
	}

	// Assuming repository has a method to update credit
	if err := repository.UserRepo.UpdateUserCredit(userID, amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user credit")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Payment captured and credit added successfully",
		"userID":  userID,
		"amount":  amount,
	})
}
