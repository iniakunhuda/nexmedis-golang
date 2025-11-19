package model

// RegisterRequest represents the request body for client registration
type RegisterRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email"`
}

// LogRequest represents the request body for logging API hits
type LogRequest struct {
	APIKey   string `json:"api_key" validate:"required"`
	IP       string `json:"ip" validate:"required,ip"`
	Endpoint string `json:"endpoint" validate:"required"`
}

// LoginRequest represents the request body for authentication
type LoginRequest struct {
	APIKey string `json:"api_key" validate:"required"`
}
