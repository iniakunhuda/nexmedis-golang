package router

import (
	"nexmedis-golang/handler"
	"nexmedis-golang/store"
	"nexmedis-golang/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// Config holds router configuration
type Config struct {
	DB                *gorm.DB
	RateLimiter       *utils.RateLimiter
	CacheTTL          time.Duration
	EnableIPWhitelist bool
	AllowedIPs        []string
}

// Setup configures all routes and middleware
func Setup(e *echo.Echo, config Config) {
	// Initialize stores
	clientStore := store.NewClientStore(config.DB)
	logStore := store.NewLogStore(config.DB)

	// Initialize handlers
	clientHandler := handler.NewClientHandler(clientStore)
	authHandler := handler.NewAuthHandler(clientStore)
	logHandler := handler.NewLogHandler(logStore, clientStore, config.RateLimiter)
	usageHandler := handler.NewUsageHandler(logStore, clientStore, config.CacheTTL)

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
	}))

	// Rate limiting middleware
	if config.RateLimiter != nil {
		e.Use(RateLimitHeaders(config.RateLimiter))
	}

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
	api.POST("/register", clientHandler.Register)
	api.POST("/login", authHandler.Login)

	// API log routes (API key required)
	api.POST("/logs", logHandler.RecordLog)

	// Protected routes (JWT required)
	protected := api.Group("")
	protected.Use(JWTMiddleware())

	// Auth routes
	protected.POST("/auth/refresh", authHandler.RefreshToken)
	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/auth/profile", authHandler.GetProfile)

	// Usage routes (JWT required)
	usage := protected.Group("/usage")

	if config.EnableIPWhitelist && len(config.AllowedIPs) > 0 {
		usage.Use(IPWhitelistMiddleware(config.AllowedIPs))
	}

	usage.GET("/daily", usageHandler.GetDailyUsage)
	usage.GET("/top", usageHandler.GetTopClients)
	usage.GET("/stats", usageHandler.GetUsageStats)
	usage.GET("/client/:client_id", usageHandler.GetClientUsage)

	// Custom error handler
	e.HTTPErrorHandler = ErrorHandler
}
