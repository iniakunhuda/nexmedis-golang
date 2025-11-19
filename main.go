package main

import (
	"context"
	"log"
	"net/http"
	"nexmedis-golang/db"
	"nexmedis-golang/router"
	"nexmedis-golang/utils"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize JWT
	utils.InitJWT()

	// Initialize database
	dbConfig := db.GetDBConfig()
	if err := db.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	redisConfig := db.GetRedisConfig()
	if err := db.InitRedis(redisConfig); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
		log.Println("Continuing without Redis cache...")
	} else {
		defer db.CloseRedis()
	}

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Get configuration
	cacheTTL := getCacheTTL()
	rateLimitPerHour := getRateLimit()

	// Initialize rate limiter
	rateLimiter := utils.NewRateLimiter(rateLimitPerHour)

	// Setup routes
	routerConfig := router.Config{
		DB:                db.DB,
		RateLimiter:       rateLimiter,
		CacheTTL:          cacheTTL,
		EnableIPWhitelist: false, // Set to true and configure AllowedIPs for IP whitelisting
		AllowedIPs:        []string{},
	}
	router.Setup(e, routerConfig)

	// Get server configuration
	serverHost := getEnv("SERVER_HOST", "0.0.0.0")
	serverPort := getEnv("SERVER_PORT", "8080")
	serverAddr := serverHost + ":" + serverPort

	// Start server with graceful shutdown
	go func() {
		log.Printf("Starting server on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getCacheTTL gets cache TTL from environment
func getCacheTTL() time.Duration {
	ttlStr := getEnv("CACHE_TTL", "3600")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return time.Hour
	}
	return time.Duration(ttl) * time.Second
}

// getRateLimit gets rate limit from environment
func getRateLimit() int {
	limitStr := getEnv("RATE_LIMIT_PER_HOUR", "1000")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 1000
	}
	return limit
}
