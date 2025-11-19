package handler

import (
	"nexmedis-golang/model"
	"nexmedis-golang/store"
	"nexmedis-golang/utils"

	"github.com/labstack/echo/v4"
)

// ClientHandler handles client-related requests
type ClientHandler struct {
	clientStore *store.ClientStore
}

// NewClientHandler creates a new ClientHandler
func NewClientHandler(clientStore *store.ClientStore) *ClientHandler {
	return &ClientHandler{
		clientStore: clientStore,
	}
}

// Register handles client registration
func (h *ClientHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	if err := utils.ValidateRequired(req.Name, "name"); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := utils.ValidateMinLength(req.Name, "name", 3); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := utils.ValidateMaxLength(req.Name, "name", 100); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Check if email already exists
	exists, err := h.clientStore.ExistsByEmail(req.Email)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to check email", err.Error())
	}
	if exists {
		return utils.ConflictResponse(c, "Email already registered")
	}

	// Generate API key
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate API key", err.Error())
	}

	// Create client
	client := &model.Client{
		Name:   utils.SanitizeString(req.Name),
		Email:  utils.SanitizeString(req.Email),
		APIKey: apiKey,
	}

	if err := h.clientStore.Create(client); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create client", err.Error())
	}

	// Return response with API key
	return utils.CreatedResponse(c, "Client registered successfully", client.ToResponse(true))
}
