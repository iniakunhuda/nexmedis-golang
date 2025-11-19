package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"nexmedis-golang/db"

	"github.com/google/uuid"
)

// RateLimiter handles rate limiting for API clients
type RateLimiter struct {
	maxRequestsPerHour int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequestsPerHour int) *RateLimiter {
	return &RateLimiter{
		maxRequestsPerHour: maxRequestsPerHour,
	}
}

// CheckLimit checks if a client has exceeded their rate limit
func (rl *RateLimiter) CheckLimit(ctx context.Context, clientID uuid.UUID) (bool, int, error) {
	key := rl.getRateLimitKey(clientID)

	// Try to get current count from Redis
	if db.IsRedisAvailable(ctx) {
		count, err := db.GetCounter(ctx, key)
		if err != nil {
			// If error, fallback to allowing the request
			return true, 0, nil
		}

		if count >= int64(rl.maxRequestsPerHour) {
			remaining := 0
			return false, remaining, nil
		}

		// Increment counter
		newCount, err := db.IncrementCounter(ctx, key)
		if err != nil {
			return true, 0, nil
		}

		// Set expiry if this is the first request
		if newCount == 1 {
			ttl := time.Hour
			_ = db.RedisClient.Expire(ctx, key, ttl).Err()
		}

		remaining := rl.maxRequestsPerHour - int(newCount)
		if remaining < 0 {
			remaining = 0
		}

		return true, remaining, nil
	}

	// If Redis is not available, allow the request
	return true, rl.maxRequestsPerHour, nil
}

// GetRemainingRequests gets the remaining requests for a client
func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, clientID uuid.UUID) (int, error) {
	key := rl.getRateLimitKey(clientID)

	if !db.IsRedisAvailable(ctx) {
		return rl.maxRequestsPerHour, nil
	}

	count, err := db.GetCounter(ctx, key)
	if err != nil {
		return rl.maxRequestsPerHour, nil
	}

	remaining := rl.maxRequestsPerHour - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// ResetLimit resets the rate limit for a client
func (rl *RateLimiter) ResetLimit(ctx context.Context, clientID uuid.UUID) error {
	key := rl.getRateLimitKey(clientID)
	return db.CacheDelete(ctx, key)
}

// getRateLimitKey generates the Redis key for rate limiting
func (rl *RateLimiter) getRateLimitKey(clientID uuid.UUID) string {
	hour := time.Now().UTC().Format("2006-01-02-15")
	return fmt.Sprintf("rate_limit:%s:%s", clientID.String(), hour)
}

// GetDefaultRateLimit returns the default rate limit from environment
func GetDefaultRateLimit() int {
	limitStr := getEnvOrDefault("RATE_LIMIT_PER_HOUR", "1000")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 1000
	}
	return limit
}

func getEnvOrDefault(key, defaultValue string) string {
	// This is a placeholder - in real implementation, use os.Getenv
	return defaultValue
}
