package handler

import (
	"encoding/json"
	"fmt"
	"nexmedis-golang/db"
	"nexmedis-golang/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// SSEHandler handles Server-Sent Events for real-time updates
type SSEHandler struct{}

// NewSSEHandler creates a new SSEHandler
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

// StreamUsageUpdates streams real-time API log updates via SSE
//
//	@Summary		Stream real-time usage updates
//	@Description	Subscribe to real-time API activity updates using Server-Sent Events (SSE). Receives notifications when new API logs are recorded.
//	@Tags			Real-time
//	@Produce		text/event-stream
//	@Security		BearerAuth
//	@Success		200	{string}	string	"Event stream connection established"
//	@Failure		401	{object}	object{success=bool,message=string,error=string}	"Unauthorized - JWT token required"
//	@Failure		503	{object}	object{success=bool,message=string,error=string}	"Redis not available for SSE"
//	@Router			/api/stream/usage [get]
func (h *SSEHandler) StreamUsageUpdates(c echo.Context) error {
	ctx := c.Request().Context()

	// Check if Redis is available
	if !db.IsRedisAvailable(ctx) {
		return utils.ServiceUnavailableResponse(c, "Redis not available for real-time updates")
	}

	// Set headers for SSE
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no") // Disable buffering in nginx

	// Get client info from context (optional)
	clientID, _ := c.Get("client_id").(string)

	// Subscribe to Redis Pub/Sub channel
	channel := "api_logs:updates"
	pubsub := db.RedisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Send initial connection message
	if err := sendSSEMessage(c, "connected", map[string]interface{}{
		"message":   "Successfully connected to real-time updates",
		"client_id": clientID,
		"timestamp": time.Now().UTC(),
	}); err != nil {
		return err
	}

	// Flush to ensure client receives the connection message
	c.Response().Flush()

	// Create a ticker for keep-alive heartbeat (every 30 seconds)
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	// Listen for messages
	msgChan := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			log.Printf("SSE client disconnected: %s", clientID)
			return nil

		case msg := <-msgChan:
			// Parse the message
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				log.Printf("Failed to parse SSE message: %v", err)
				continue
			}

			// Send update to client
			if err := sendSSEMessage(c, "update", data); err != nil {
				log.Printf("Failed to send SSE message: %v", err)
				return err
			}

			c.Response().Flush()

		case <-heartbeat.C:
			// Send heartbeat to keep connection alive
			if err := sendSSEMessage(c, "heartbeat", map[string]interface{}{
				"timestamp": time.Now().UTC(),
			}); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
				return err
			}

			c.Response().Flush()
		}
	}
}

// StreamTopClients streams real-time top clients updates via SSE
//
//	@Summary		Stream top clients updates
//	@Description	Subscribe to real-time updates of top clients with highest API usage. Updates are sent every 60 seconds.
//	@Tags			Real-time
//	@Produce		text/event-stream
//	@Security		BearerAuth
//	@Success		200	{string}	string	"Event stream connection established"
//	@Failure		401	{object}	object{success=bool,message=string,error=string}	"Unauthorized - JWT token required"
//	@Router			/api/stream/top [get]
func (h *SSEHandler) StreamTopClients(c echo.Context) error {
	ctx := c.Request().Context()

	// Set headers for SSE
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	// Get client info from context
	clientID, _ := c.Get("client_id").(string)

	// Send initial connection message
	if err := sendSSEMessage(c, "connected", map[string]interface{}{
		"message":   "Successfully connected to top clients stream",
		"client_id": clientID,
		"timestamp": time.Now().UTC(),
	}); err != nil {
		return err
	}

	c.Response().Flush()

	// Create ticker for periodic updates (every 60 seconds)
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Create heartbeat ticker (every 30 seconds)
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			log.Printf("SSE client disconnected: %s", clientID)
			return nil

		case <-ticker.C:
			// Check if cache has top clients data
			if db.IsRedisAvailable(ctx) {
				var topClients []map[string]interface{}
				if err := db.CacheGet(ctx, "usage:top:24h", &topClients); err == nil {
					// Send cached top clients
					if err := sendSSEMessage(c, "top_clients", map[string]interface{}{
						"data":      topClients,
						"timestamp": time.Now().UTC(),
						"source":    "cache",
					}); err != nil {
						log.Printf("Failed to send top clients update: %v", err)
						return err
					}
				} else {
					// Send notification that data is being refreshed
					if err := sendSSEMessage(c, "refreshing", map[string]interface{}{
						"message":   "Top clients data is being refreshed",
						"timestamp": time.Now().UTC(),
					}); err != nil {
						return err
					}
				}
			}

			c.Response().Flush()

		case <-heartbeat.C:
			// Send heartbeat
			if err := sendSSEMessage(c, "heartbeat", map[string]interface{}{
				"timestamp": time.Now().UTC(),
			}); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
				return err
			}

			c.Response().Flush()
		}
	}
}


// sendSSEMessage sends an SSE formatted message
func sendSSEMessage(c echo.Context, event string, data interface{}) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write SSE formatted message
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event, string(jsonData))
	_, err = c.Response().Write([]byte(message))
	return err
}
