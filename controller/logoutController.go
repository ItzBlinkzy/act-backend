package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = os.Getenv("AUTH_COOKIE")
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0) // set the expiration to a past time
	cookie.HttpOnly = true
	cookie.Secure = false // Should be true over HTTPS
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, "User logged out successfully")
}
