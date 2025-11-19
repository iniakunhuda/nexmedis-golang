package handler

import (
	"context"
	"nexmedis-golang/db"
	"nexmedis-golang/model"
	"nexmedis-golang/store"
	"nexmedis-golang/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// LogHandler handles API log-related requests
type LogHandler struct {
	logStore    *store.LogStore
	clientStore *store.ClientStore
	rateLimiter *utils.RateLimiter
}

// NewLogHandler creates a new LogHandler
func NewLogHandler(logStore *store.LogStore, clientStore *store.ClientStore, rateLimiter *utils.RateLimiter) *LogHandler {
	return &LogHandler{
		logStore:    logStore,
		clientStore: clientStore,
		rateLimiter: rateLimiter,
	}
}

// RecordLog handles recording an API hit
func (h *LogHandler) RecordLog(c echo.Context) error {
	var req model.LogRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	if err := utils.ValidateAPIKey(req.APIKey); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := utils.ValidateIP(req.IP); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := utils.ValidateEndpoint(req.Endpoint); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Find client by API key
	client, err := h.clientStore.FindByAPIKey(req.APIKey)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid API key")
	}

	// Check rate limit
	ctx := c.Request().Context()
	allowed, remaining, err := h.rateLimiter.CheckLimit(ctx, client.ID)
	if err != nil {
		// Log error but continue (graceful degradation)
	}

	if !allowed {
		return utils.TooManyRequestsResponse(c, "Rate limit exceeded")
	}

	// Create log entry
	log := &model.APILog{
		ClientID:  client.ID,
		APIKey:    req.APIKey,
		IP:        req.IP,
		Endpoint:  req.Endpoint,
		Timestamp: time.Now().UTC(),
	}

	if err := h.logStore.Create(log); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to record log", err.Error())
	}

	// Invalidate cache for usage endpoints
	go h.invalidateUsageCache(ctx, client.ID)

	// Publish update via Redis Pub/Sub
	go h.publishLogUpdate(ctx, log)

	response := map[string]interface{}{
		"log_id":             log.ID,
		"timestamp":          log.Timestamp,
		"remaining_requests": remaining,
	}

	return utils.CreatedResponse(c, "API hit recorded successfully", response)
}

// invalidateUsageCache invalidates usage-related cache entries
func (h *LogHandler) invalidateUsageCache(ctx context.Context, clientID interface{}) {
	bgCtx := context.Background()

	if !db.IsRedisAvailable(bgCtx) {
		log.Warn("Redis not available for cache invalidation")
		return
	}

	if err := db.CacheInvalidatePattern(bgCtx, "usage:daily:*"); err != nil {
		log.Printf("Failed to invalidate daily usage cache: %v", err)
	}

	if err := db.CacheDelete(bgCtx, "usage:top:24h"); err != nil {
		log.Printf("Failed to invalidate top clients cache: %v", err)
	}

	log.Printf("Cache invalidated for client: %v", clientID)
}

// publishLogUpdate publishes a log update to Redis Pub/Sub
func (h *LogHandler) publishLogUpdate(ctx context.Context, log *model.APILog) {
	if !db.IsRedisAvailable(ctx) {
		return
	}

	message := map[string]interface{}{
		"client_id": log.ClientID,
		"endpoint":  log.Endpoint,
		"timestamp": log.Timestamp,
	}

	_ = db.PublishMessage(ctx, "api_logs:updates", message)
}
