package handler

import (
	"context"
	"nexmedis-golang/db"
	"nexmedis-golang/model"
	"nexmedis-golang/store"
	"nexmedis-golang/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// UsageHandler handles usage statistics requests
type UsageHandler struct {
	logStore    *store.LogStore
	clientStore *store.ClientStore
	cacheTTL    time.Duration
}

// NewUsageHandler creates a new UsageHandler
func NewUsageHandler(logStore *store.LogStore, clientStore *store.ClientStore, cacheTTL time.Duration) *UsageHandler {
	return &UsageHandler{
		logStore:    logStore,
		clientStore: clientStore,
		cacheTTL:    cacheTTL,
	}
}

// GetDailyUsage returns daily usage statistics for the last 7 days
func (h *UsageHandler) GetDailyUsage(c echo.Context) error {
	ctx := c.Request().Context()
	cacheKey := "usage:daily:7days"

	// Try to get from cache
	if db.IsRedisAvailable(ctx) {
		var cachedData []model.DailyUsage
		if err := db.CacheGet(ctx, cacheKey, &cachedData); err == nil {
			return utils.OKResponse(c, "Daily usage retrieved from cache", cachedData)
		}
	}

	// Get from database
	usage, err := h.logStore.GetDailyUsage(7)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get daily usage", err.Error())
	}

	// Cache the result
	if db.IsRedisAvailable(ctx) {
		_ = db.CacheSet(ctx, cacheKey, usage, h.cacheTTL)
	}

	return utils.OKResponse(c, "Daily usage retrieved successfully", usage)
}

// GetTopClients returns top 3 clients with the highest requests in the last 24 hours
func (h *UsageHandler) GetTopClients(c echo.Context) error {
	ctx := c.Request().Context()
	cacheKey := "usage:top:24h"

	// Try to get from cache
	if db.IsRedisAvailable(ctx) {
		var cachedData []model.TopClient
		if err := db.CacheGet(ctx, cacheKey, &cachedData); err == nil {
			return utils.OKResponse(c, "Top clients retrieved from cache", cachedData)
		}
	}

	// Prefetch mechanism: check if cache is about to expire and refresh it
	go h.prefetchTopClients(ctx, cacheKey)

	// Get from database
	topClients, err := h.logStore.GetTopClients(3, 24*time.Hour)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get top clients", err.Error())
	}

	// Cache the result
	if db.IsRedisAvailable(ctx) {
		_ = db.CacheSet(ctx, cacheKey, topClients, h.cacheTTL)
	}

	return utils.OKResponse(c, "Top clients retrieved successfully", topClients)
}

// prefetchTopClients implements cache prefetching to avoid cache misses
func (h *UsageHandler) prefetchTopClients(ctx context.Context, cacheKey string) {
	if !db.IsRedisAvailable(ctx) {
		return
	}

	// Check TTL of the cache key
	ttl, err := db.RedisClient.TTL(ctx, cacheKey).Result()
	if err != nil {
		return
	}

	// If TTL is less than 5 minutes, refresh the cache
	if ttl > 0 && ttl < 5*time.Minute {
		topClients, err := h.logStore.GetTopClients(3, 24*time.Hour)
		if err == nil {
			_ = db.CacheSet(ctx, cacheKey, topClients, h.cacheTTL)
		}
	}
}

// GetClientUsage returns usage statistics for a specific client
func (h *UsageHandler) GetClientUsage(c echo.Context) error {
	clientIDStr := c.Param("client_id")
	if clientIDStr == "" {
		return utils.BadRequestResponse(c, "Client ID is required")
	}

	// Find client
	client, err := h.clientStore.FindByClientID(clientIDStr)
	if err != nil {
		return utils.NotFoundResponse(c, "Client not found")
	}

	ctx := c.Request().Context()
	cacheKey := "usage:client:" + clientIDStr + ":7days"

	// Try to get from cache
	if db.IsRedisAvailable(ctx) {
		var cachedData []model.DailyUsage
		if err := db.CacheGet(ctx, cacheKey, &cachedData); err == nil {
			return utils.OKResponse(c, "Client usage retrieved from cache", cachedData)
		}
	}

	// Get from database
	usage, err := h.logStore.GetDailyUsageByClient(client.ID, 7)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get client usage", err.Error())
	}

	// Cache the result
	if db.IsRedisAvailable(ctx) {
		_ = db.CacheSet(ctx, cacheKey, usage, h.cacheTTL)
	}

	return utils.OKResponse(c, "Client usage retrieved successfully", usage)
}

// GetUsageStats returns overall usage statistics
func (h *UsageHandler) GetUsageStats(c echo.Context) error {
	ctx := c.Request().Context()
	cacheKey := "usage:stats:overall"

	// Try to get from cache
	if db.IsRedisAvailable(ctx) {
		var cachedData map[string]interface{}
		if err := db.CacheGet(ctx, cacheKey, &cachedData); err == nil {
			return utils.OKResponse(c, "Usage stats retrieved from cache", cachedData)
		}
	}

	now := time.Now().UTC()
	last24h := now.Add(-24 * time.Hour)
	last7d := now.Add(-7 * 24 * time.Hour)

	// Get counts
	count24h, err := h.logStore.GetTotalRequestCount(last24h, now)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get 24h count", err.Error())
	}

	count7d, err := h.logStore.GetTotalRequestCount(last7d, now)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get 7d count", err.Error())
	}

	// Get total clients
	totalClients, err := h.clientStore.Count()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get client count", err.Error())
	}

	stats := map[string]interface{}{
		"total_requests_24h": count24h,
		"total_requests_7d":  count7d,
		"total_clients":      totalClients,
		"timestamp":          now,
	}

	// Cache the result
	if db.IsRedisAvailable(ctx) {
		_ = db.CacheSet(ctx, cacheKey, stats, 5*time.Minute) // Shorter TTL for stats
	}

	return utils.OKResponse(c, "Usage stats retrieved successfully", stats)
}
