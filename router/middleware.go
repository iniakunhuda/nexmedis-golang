package router

import (
	"nexmedis-golang/utils"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates JWT tokens
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.UnauthorizedResponse(c, "Authorization header required")
			}

			// Extract token
			tokenString := ""
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				return utils.UnauthorizedResponse(c, "Invalid authorization header format")
			}

			// Validate token
			claims, err := utils.ValidateJWT(tokenString)
			if err != nil {
				return utils.UnauthorizedResponse(c, "Invalid or expired token")
			}

			// Set client information in context
			c.Set("client_id", claims.ClientID.String())
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

func RateLimitHeaders(limiter *utils.RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Execute the handler first
			err := next(c)

			// Add rate limit headers if client_id is in context
			if clientIDStr, ok := c.Get("client_id").(string); ok {
				// Parse client ID and get remaining requests
				// This is a simplified version
				c.Response().Header().Set("X-RateLimit-Limit", "1000")
				c.Response().Header().Set("X-RateLimit-Remaining", "999")
				_ = clientIDStr // Use the variable to avoid unused error
			}

			return err
		}
	}
}

func ErrorHandler(err error, c echo.Context) {
	code := 500
	message := "Internal server error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Don't send error response if response already started
	if !c.Response().Committed {
		_ = utils.ErrorResponseWithMessage(c, code, message, "")
	}
}