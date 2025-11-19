package model

// RegisterRequest represents the request body for client registration
// @Description Request body for registering a new client
type RegisterRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100" example:"John Doe"`      // Client name (3-100 characters)
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"` // Valid email address
}

// LogRequest represents the request body for logging API hits
// @Description Request body for recording an API hit
type LogRequest struct {
	APIKey   string `json:"api_key" validate:"required" example:"sk_live_abcdef123456"` // Client's API key
	IP       string `json:"ip" validate:"required,ip" example:"192.168.1.100"`          // Client's IP address
	Endpoint string `json:"endpoint" validate:"required" example:"/api/v1/users"`       // API endpoint that was called
}

// LoginRequest represents the request body for authentication
// @Description Request body for client authentication
type LoginRequest struct {
	APIKey string `json:"api_key" validate:"required" example:"sk_live_abcdef123456"` // API key for authentication
}
