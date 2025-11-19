package db

import (
	"fmt"
	"log"
	"nexmedis-golang/model"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// InitDB initializes the database connection
func InitDB(config Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")

	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&model.Client{},
		&model.APILog{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create indexes for better query performance
	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createIndexes creates additional database indexes
func createIndexes() error {
	// Index for api_logs timestamp queries
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_api_logs_timestamp 
		ON api_logs(timestamp DESC)
	`).Error; err != nil {
		return err
	}

	// Composite index for client_id and timestamp (for daily usage queries)
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_api_logs_client_timestamp 
		ON api_logs(client_id, timestamp DESC)
	`).Error; err != nil {
		return err
	}

	// Index for endpoint queries
	if err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_api_logs_endpoint 
		ON api_logs(endpoint)
	`).Error; err != nil {
		return err
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// GetDBConfig loads database configuration from environment variables
func GetDBConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "api_tracking"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
