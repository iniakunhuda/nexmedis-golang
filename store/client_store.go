package store

import (
	"errors"
	"nexmedis-golang/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClientStore handles database operations for clients
type ClientStore struct {
	db *gorm.DB
}

// NewClientStore creates a new ClientStore instance
func NewClientStore(db *gorm.DB) *ClientStore {
	return &ClientStore{db: db}
}

// Create creates a new client in the database
func (s *ClientStore) Create(client *model.Client) error {
	return s.db.Create(client).Error
}

// FindByID finds a client by UUID
func (s *ClientStore) FindByID(id uuid.UUID) (*model.Client, error) {
	var client model.Client
	err := s.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &client, nil
}

// FindByClientID finds a client by client_id string
func (s *ClientStore) FindByClientID(clientID string) (*model.Client, error) {
	var client model.Client
	err := s.db.Where("client_id = ?", clientID).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &client, nil
}

// FindByAPIKey finds a client by API key
func (s *ClientStore) FindByAPIKey(apiKey string) (*model.Client, error) {
	var client model.Client
	err := s.db.Where("api_key = ?", apiKey).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &client, nil
}

// FindByEmail finds a client by email
func (s *ClientStore) FindByEmail(email string) (*model.Client, error) {
	var client model.Client
	err := s.db.Where("email = ?", email).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &client, nil
}

// ExistsByEmail checks if a client with the given email exists
func (s *ClientStore) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := s.db.Model(&model.Client{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByAPIKey checks if a client with the given API key exists
func (s *ClientStore) ExistsByAPIKey(apiKey string) (bool, error) {
	var count int64
	err := s.db.Model(&model.Client{}).Where("api_key = ?", apiKey).Count(&count).Error
	return count > 0, err
}

// List returns all clients with pagination
func (s *ClientStore) List(offset, limit int) ([]model.Client, error) {
	var clients []model.Client
	err := s.db.Offset(offset).Limit(limit).Find(&clients).Error
	return clients, err
}

// Update updates a client
func (s *ClientStore) Update(client *model.Client) error {
	return s.db.Save(client).Error
}

// Delete soft deletes a client
func (s *ClientStore) Delete(id uuid.UUID) error {
	return s.db.Delete(&model.Client{}, id).Error
}

// Count returns the total number of clients
func (s *ClientStore) Count() (int64, error) {
	var count int64
	err := s.db.Model(&model.Client{}).Count(&count).Error
	return count, err
}
