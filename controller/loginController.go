// controller/login_controller.go

package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	fmt.Println("%+v", c)

	// Log request headers for debugging purposes
	fmt.Println("Request Headers:", c.Request().Header)

	// Log request IP address (origin)
	ip := c.RealIP()
	fmt.Println("Client IP:", ip)

	// Log user agent
	userAgent := c.Request().Header.Get("User-Agent")
	fmt.Println("User-Agent:", userAgent)
	var payload model.LoginPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	user, err := repository.UserRepo.FindByEmail(strings.ToLower(payload.Email))
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	if (user.LoginMethod == "oauth") {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Please login using Google or GitHub"})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}
	// generating token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Message: "An error occurred. Please try again or contact the administrator."})
	}

	associations, err := repository.GetAllClients(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not retrieve client-manager associations"})
	}

	
	cookie := &http.Cookie{}
	cookie.Name = os.Getenv("AUTH_COOKIE")
	cookie.Value = t
	cookie.Expires = time.Now().Add(3 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	// Return JSON response with token and associations
	response := model.LoginResponse{
		Token:                  t,
		Clients: associations,
	}
	return c.JSON(http.StatusOK, response)

}
