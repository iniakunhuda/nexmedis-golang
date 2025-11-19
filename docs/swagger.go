// Package docs provides API documentation configuration
//
//	@title						User Activity Tracking System API
//	@version					1.0
//	@description				A high-performance API for tracking user activity with advanced caching, rate limiting, and JWT authentication. Built with Go, Echo, PostgreSQL, and Redis.
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.url				http://www.nexmedis.com/support
//	@contact.email				support@nexmedis.com
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
//	@description				API Key for client authentication
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token
//
//	@schemes					http https
//
//	@x-extension-openapi		{"example": "value on a json format"}
package docs
