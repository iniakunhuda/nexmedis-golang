package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Client represents a registered API client
type Client struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ClientID  string         `gorm:"uniqueIndex;not null" json:"client_id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	APIKey    string         `gorm:"uniqueIndex;not null" json:"-"` // Don't expose in JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID and ClientID
func (c *Client) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ClientID == "" {
		c.ClientID = "client_" + uuid.New().String()[:8]
	}
	return nil
}

// TableName specifies the table name for Client
func (Client) TableName() string {
	return "clients"
}

// ClientResponse is used for API responses (excludes sensitive data)
// @Description Client information response
type ClientResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` // Client UUID
	ClientID  string    `json:"client_id" example:"client_abc12345"`               // Human-readable client ID
	Name      string    `json:"name" example:"John Doe"`                           // Client name
	Email     string    `json:"email" example:"john.doe@example.com"`              // Client email
	APIKey    string    `json:"api_key,omitempty" example:"sk_live_abcdef123456"`  // API key (only shown on registration)
	CreatedAt time.Time `json:"created_at" example:"2025-01-15T10:30:00Z"`         // Creation timestamp
}

// ToResponse converts Client to ClientResponse
func (c *Client) ToResponse(includeAPIKey bool) *ClientResponse {
	resp := &ClientResponse{
		ID:        c.ID,
		ClientID:  c.ClientID,
		Name:      c.Name,
		Email:     c.Email,
		CreatedAt: c.CreatedAt,
	}
	if includeAPIKey {
		resp.APIKey = c.APIKey
	}
	return resp
}
