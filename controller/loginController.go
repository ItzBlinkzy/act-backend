// controller/login_controller.go

package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	fmt.Println("%+v", c)
	fmt.Println("Request Headers:", c.Request().Header)

	ip := c.RealIP()
	fmt.Println("Client IP:", ip)

	// Log user agent
	userAgent := c.Request().Header.Get("User-Agent")
	fmt.Println("User-Agent:", userAgent)
	var payload model.LoginPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	user, err := repository.UserRepo.FindByEmail(payload.Email)
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
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

	cookie := &http.Cookie{}
	cookie.Name = os.Getenv("AUTH_COOKIE")
	cookie.Value = t
	expirationTime := time.Now().Add(10 * time.Hour) // Matching cookie expiration
	claims["exp"] = expirationTime.Unix()
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, nil)

}
