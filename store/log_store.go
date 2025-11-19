package store

import (
	"nexmedis-golang/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogStore handles database operations for API logs
type LogStore struct {
	db *gorm.DB
}

// NewLogStore creates a new LogStore instance
func NewLogStore(db *gorm.DB) *LogStore {
	return &LogStore{db: db}
}

// Create creates a new API log entry
func (s *LogStore) Create(log *model.APILog) error {
	return s.db.Create(log).Error
}

// BatchCreate creates multiple API log entries in a single transaction
func (s *LogStore) BatchCreate(logs []model.APILog) error {
	if len(logs) == 0 {
		return nil
	}
	return s.db.CreateInBatches(logs, 100).Error
}

// FindByID finds an API log by ID
func (s *LogStore) FindByID(id uuid.UUID) (*model.APILog, error) {
	var log model.APILog
	err := s.db.Where("id = ?", id).First(&log).Error
	return &log, err
}

// GetDailyUsage returns daily usage for each client for the last N days
func (s *LogStore) GetDailyUsage(days int) ([]model.DailyUsage, error) {
	var results []model.DailyUsage

	startDate := time.Now().UTC().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	query := `
		SELECT 
			l.client_id,
			c.name as client_name,
			DATE(l.timestamp) as date,
			COUNT(*) as count
		FROM api_logs l
		INNER JOIN clients c ON l.client_id = c.id
		WHERE l.timestamp >= ?
		GROUP BY l.client_id, c.name, DATE(l.timestamp)
		ORDER BY date DESC, count DESC
	`

	err := s.db.Raw(query, startDate).Scan(&results).Error
	return results, err
}

// GetDailyUsageByClient returns daily usage for a specific client
func (s *LogStore) GetDailyUsageByClient(clientID uuid.UUID, days int) ([]model.DailyUsage, error) {
	var results []model.DailyUsage

	startDate := time.Now().UTC().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	query := `
		SELECT 
			l.client_id,
			c.name as client_name,
			DATE(l.timestamp) as date,
			COUNT(*) as count
		FROM api_logs l
		INNER JOIN clients c ON l.client_id = c.id
		WHERE l.client_id = ? AND l.timestamp >= ?
		GROUP BY l.client_id, c.name, DATE(l.timestamp)
		ORDER BY date DESC
	`

	err := s.db.Raw(query, clientID, startDate).Scan(&results).Error
	return results, err
}

// GetTopClients returns top N clients by request count in the last duration
func (s *LogStore) GetTopClients(limit int, duration time.Duration) ([]model.TopClient, error) {
	var results []model.TopClient

	startTime := time.Now().UTC().Add(-duration)

	query := `
		SELECT 
			l.client_id,
			c.name as client_name,
			c.email,
			COUNT(*) as total_requests
		FROM api_logs l
		INNER JOIN clients c ON l.client_id = c.id
		WHERE l.timestamp >= ?
		GROUP BY l.client_id, c.name, c.email
		ORDER BY total_requests DESC
		LIMIT ?
	`

	err := s.db.Raw(query, startTime, limit).Scan(&results).Error
	return results, err
}

// GetClientRequestCount returns the total request count for a client in a time range
func (s *LogStore) GetClientRequestCount(clientID uuid.UUID, start, end time.Time) (int64, error) {
	var count int64
	err := s.db.Model(&model.APILog{}).
		Where("client_id = ? AND timestamp >= ? AND timestamp < ?", clientID, start, end).
		Count(&count).Error
	return count, err
}

// GetTotalRequestCount returns total requests in a time range
func (s *LogStore) GetTotalRequestCount(start, end time.Time) (int64, error) {
	var count int64
	err := s.db.Model(&model.APILog{}).
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Count(&count).Error
	return count, err
}

// GetRequestsByEndpoint returns request count grouped by endpoint
func (s *LogStore) GetRequestsByEndpoint(start, end time.Time) (map[string]int64, error) {
	type Result struct {
		Endpoint string
		Count    int64
	}

	var results []Result
	err := s.db.Model(&model.APILog{}).
		Select("endpoint, COUNT(*) as count").
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Group("endpoint").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	endpointCounts := make(map[string]int64)
	for _, r := range results {
		endpointCounts[r.Endpoint] = r.Count
	}

	return endpointCounts, nil
}

// DeleteOldLogs deletes logs older than the specified duration
func (s *LogStore) DeleteOldLogs(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().UTC().Add(-olderThan)
	result := s.db.Where("timestamp < ?", cutoffTime).Delete(&model.APILog{})
	return result.RowsAffected, result.Error
}

// List returns API logs with pagination
func (s *LogStore) List(offset, limit int) ([]model.APILog, error) {
	var logs []model.APILog
	err := s.db.Offset(offset).Limit(limit).Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

// ListByClient returns API logs for a specific client with pagination
func (s *LogStore) ListByClient(clientID uuid.UUID, offset, limit int) ([]model.APILog, error) {
	var logs []model.APILog
	err := s.db.Where("client_id = ?", clientID).
		Offset(offset).
		Limit(limit).
		Order("timestamp DESC").
		Find(&logs).Error
	return logs, err
}
