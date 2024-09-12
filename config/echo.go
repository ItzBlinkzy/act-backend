package config

import (
	"log"
	"net/http"
	"os"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func isSkippedPath(c echo.Context) bool {
	skippedPaths := []string{
		"/api/v1/login",
		"/api/v1/register-company-first-user",
	}
	for _, path := range skippedPaths {
		if path == c.Path() {
			return true
		}
	}

	if os.Getenv("ENV") == "development" {
		return true
	}

	return false
}

func refreshToken(c echo.Context) {
	cookie, err := c.Cookie(os.Getenv("AUTH_COOKIE"))
	if err == nil && cookie != nil {
		// Extend its expiration by 3 hours
		cookie.Expires = time.Now().Add(3 * time.Hour)
		cookie.Path = "/"
		c.SetCookie(cookie) // Reset the cookie with the new expiration
	}
}

func InitEcho() *echo.Echo {
	e := echo.New()

	var allowOrigins []string
	var allowHeaders []string

	if os.Getenv("ENV") == "development" {
		allowOrigins = []string{"http://localhost:5173", "http://192.168.0.67:5173", "http://192.168.0.210:5173"}
		allowHeaders = []string{"Content-Type", "Timezone", "User-email"}
	} else {
		allowOrigins = []string{"https://condominioforyou.app", "https://www.condominioforyou.app"}
		allowHeaders = []string{"Content-Type", "Timezone"}
	}

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 5, Burst: 10, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println("Request timed out:", c.Path())
		},
		Timeout: 25 * time.Second,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodHead, http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
	}))

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:     []byte(os.Getenv("JWT_SECRET")),
		TokenLookup:    "cookie:" + os.Getenv("AUTH_COOKIE"),
		Skipper:        isSkippedPath,
		SuccessHandler: refreshToken,
	}))

	return e
}
