package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APILog represents an API request log entry
type APILog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ClientID  uuid.UUID `gorm:"type:uuid;index:idx_client_timestamp;not null" json:"client_id"`
	APIKey    string    `gorm:"index;not null" json:"-"`
	IP        string    `gorm:"not null" json:"ip"`
	Endpoint  string    `gorm:"index;not null" json:"endpoint"`
	Timestamp time.Time `gorm:"index:idx_client_timestamp;not null" json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook to generate UUID and set timestamp
func (l *APILog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.Timestamp.IsZero() {
		l.Timestamp = time.Now()
	}
	return nil
}

// TableName specifies the table name for APILog
func (APILog) TableName() string {
	return "api_logs"
}

// DailyUsage represents daily usage statistics for a client
type DailyUsage struct {
	ClientID  uuid.UUID `json:"client_id"`
	ClientName string   `json:"client_name"`
	Date      string    `json:"date"`
	Count     int64     `json:"count"`
}

// TopClient represents top clients by request count
type TopClient struct {
	ClientID   uuid.UUID `json:"client_id"`
	ClientName string    `json:"client_name"`
	Email      string    `json:"email"`
	TotalRequests int64  `json:"total_requests"`
}
