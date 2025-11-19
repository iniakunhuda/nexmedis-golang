package handler

import (
	"nexmedis-golang/model"
	"nexmedis-golang/store"
	"nexmedis-golang/utils"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	clientStore *store.ClientStore
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(clientStore *store.ClientStore) *AuthHandler {
	return &AuthHandler{
		clientStore: clientStore,
	}
}

// Login handles client authentication and returns a JWT token
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate API key
	if err := utils.ValidateAPIKey(req.APIKey); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Find client by API key
	client, err := h.clientStore.FindByAPIKey(req.APIKey)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid API key")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(client.ID, client.Email)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate token", err.Error())
	}

	response := map[string]interface{}{
		"token":      token,
		"client_id":  client.ClientID,
		"expires_in": "24h",
	}

	return utils.OKResponse(c, "Login successful", response)
}

// RefreshToken handles JWT token refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// Get current token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.UnauthorizedResponse(c, "Authorization header required")
	}

	// Extract token (assuming "Bearer <token>")
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		return utils.UnauthorizedResponse(c, "Invalid authorization header format")
	}

	// Refresh token
	newToken, err := utils.RefreshJWT(tokenString)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid or expired token")
	}

	response := map[string]interface{}{
		"token":      newToken,
		"expires_in": "24h",
	}

	return utils.OKResponse(c, "Token refreshed successfully", response)
}

// Logout handles client logout (token invalidation would go here)
func (h *AuthHandler) Logout(c echo.Context) error {
	// In a production system, you would invalidate the token here
	// This could involve adding it to a blacklist in Redis
	return utils.OKResponse(c, "Logout successful", nil)
}

// GetProfile returns the authenticated client's profile
func (h *AuthHandler) GetProfile(c echo.Context) error {
	// Get client from context (set by auth middleware)
	clientID, ok := c.Get("client_id").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Client not found in context")
	}

	// Find client
	client, err := h.clientStore.FindByClientID(clientID)
	if err != nil {
		return utils.NotFoundResponse(c, "Client not found")
	}

	return utils.OKResponse(c, "Profile retrieved successfully", client.ToResponse(false))
}
