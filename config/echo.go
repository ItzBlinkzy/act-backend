package config

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Initialize OAuth2 configuration
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/api/v1/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// Handler to initiate Google login
func GoogleLoginHandler(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Handler for Google OAuth2 callback
func GoogleCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to exchange token"})
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get user info"})
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read response body"})
	}
	return c.JSON(http.StatusOK, echo.Map{"data": string(contents)})
}

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
		allowOrigins = []string{"https://condominioforyou.app", "https://www.condominioforyou.app", "http://localhost:5173", "https://act-frontend.netlify.app"}
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
