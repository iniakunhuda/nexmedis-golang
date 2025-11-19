package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global Redis client instance
var RedisClient *redis.Client

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// InitRedis initializes the Redis connection
func InitRedis(config RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection established successfully")
	return nil
}

// GetRedisConfig loads Redis configuration from environment variables
func GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	}
}

// CacheSet sets a value in Redis with TTL
func CacheSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return RedisClient.Set(ctx, key, data, ttl).Err()
}

// CacheGet retrieves a value from Redis
func CacheGet(ctx context.Context, key string, dest interface{}) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// CacheDelete deletes a key from Redis
func CacheDelete(ctx context.Context, keys ...string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return RedisClient.Del(ctx, keys...).Err()
}

// CacheInvalidatePattern invalidates all keys matching a pattern
func CacheInvalidatePattern(ctx context.Context, pattern string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// IncrementCounter increments a counter in Redis (atomic operation)
func IncrementCounter(ctx context.Context, key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return RedisClient.Incr(ctx, key).Result()
}

// GetCounter gets the current value of a counter
func GetCounter(ctx context.Context, key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	val, err := RedisClient.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// PublishMessage publishes a message to a Redis channel
func PublishMessage(ctx context.Context, channel string, message interface{}) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return RedisClient.Publish(ctx, channel, data).Err()
}

// SubscribeChannel subscribes to a Redis channel
func SubscribeChannel(ctx context.Context, channel string) *redis.PubSub {
	if RedisClient == nil {
		return nil
	}

	return RedisClient.Subscribe(ctx, channel)
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Close()
}

// IsRedisAvailable checks if Redis is available
func IsRedisAvailable(ctx context.Context) bool {
	if RedisClient == nil {
		return false
	}
	return RedisClient.Ping(ctx).Err() == nil
}
