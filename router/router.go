package router

import (
	"nexmedis-golang/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// Config holds router configuration
type Config struct {
	DB                *gorm.DB
	CacheTTL          time.Duration
	EnableIPWhitelist bool
	AllowedIPs        []string
}

func ErrorHandler(err error, c echo.Context) {
	code := 500
	message := "Internal server error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Don't send error response if response already started
	if !c.Response().Committed {
		_ = utils.ErrorResponseWithMessage(c, code, message, "")
	}
}

// Setup configures all routes and middleware
func Setup(e *echo.Echo, config Config) {

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
	}))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"time":   time.Now().UTC().String(),
		})
	})

	// API routes
	api := e.Group("/api")

	// Public routes (no authentication required)
	api.POST("/register", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	api.POST("/login", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// API log routes (API key required)
	api.POST("/logs", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	api.GET("/usage/daily", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	api.GET("/usage/top", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Custom error handler
	e.HTTPErrorHandler = ErrorHandler
}
