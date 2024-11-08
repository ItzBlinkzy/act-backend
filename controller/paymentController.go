package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"github.com/stripe/stripe-go/webhook"
)

func AddCredit(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid amount")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Account Credit"),
					},
					UnitAmount: stripe.Int64(int64(amount * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String("payment"),
		SuccessURL:        stripe.String("https://yourwebsite.com/success"),
		CancelURL:         stripe.String("https://yourwebsite.com/cancel"),
		ClientReferenceID: stripe.String(strconv.Itoa(userID)),
	}

	// Create a new checkout session
	sess, err := session.New(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create Stripe session")
	}

	// Generate the checkout URL using the session ID
	checkoutURL := "https://checkout.stripe.com/pay/" + sess.ID

	return c.JSON(http.StatusOK, map[string]interface{}{
		"checkout_url": checkoutURL,
	})
}

func UpdateCredit(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid amount")
	}

	if err := repository.UserRepo.UpdateUserCredit(userID, amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user credit")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Credit added successfully",
		"amount":  amount,
	})
}

func HandleStripeWebhook(c echo.Context) error {
	const MaxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Error reading request body")
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, c.Request().Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Webhook verification failed")
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse webhook data")
		}

		// Retrieve the user ID from ClientReferenceID
		userID, err := strconv.Atoi(session.ClientReferenceID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID in ClientReferenceID")
		}

		// Update user's credit in the database
		amount := float64(session.AmountTotal) / 100.0 // Stripe amount is in cents
		if err := repository.UserRepo.UpdateUserCredit(userID, amount); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user credit")
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
