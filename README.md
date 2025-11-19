# API Activity Tracking System

A high-performance, scalable backend system for tracking API usage with advanced caching, rate limiting, and JWT authentication built with Go, Echo, GORM, PostgreSQL, and Redis.

## 🚀 Features

### Core Features
- **RESTful API Design** - Clean, intuitive API endpoints
- **Client Management** - Register and manage API clients
- **Activity Logging** - Track all API hits with detailed metadata
- **Usage Analytics** - Daily usage reports and top client statistics
- **JWT Authentication** - Secure token-based authentication
- **API Key Management** - Generate and validate API keys

### Advanced Features
- **Redis Caching** - High-performance caching with TTL and invalidation
- **Rate Limiting** - Per-client hourly rate limits (default: 1000 req/hour)
- **Database Optimization** - Indexed queries, batch operations
- **Graceful Degradation** - Fallback when Redis is unavailable

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Redis 6+

## 🛠️ Installation

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/iniakunhuda/nexmedis-golang
cd nexmedis-golang
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run locally**
```bash
# Start PostgreSQL and Redis first
go run main.go
```

## 📡 API Endpoints

### Public Endpoints

#### Register a New Client
```http
POST /api/register
Content-Type: application/json

{
  "name": "Huda",
  "email": "huda@gmail.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Client registered successfully",
  "data": {
    "id": "uuid",
    "client_id": "client_abc12345",
    "name": "Huda",
    "email": "huda@gmail.com",
    "api_key": "generated-api-key",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "api_key": "your-api-key"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "jwt-token",
    "client_id": "client_abc12345",
    "expires_in": "24h"
  }
}
```

#### Record API Hit
```http
POST /api/logs
Content-Type: application/json

{
  "api_key": "your-api-key",
  "ip": "192.168.1.1",
  "endpoint": "/api/some-endpoint"
}
```

### Protected Endpoints (Require JWT Token)

Add JWT token to all protected endpoints:
```http
Authorization: Bearer your-jwt-token
```

#### Get Daily Usage (Last 7 Days)
```http
GET /api/usage/daily
```

**Response:**
```json
{
  "success": true,
  "message": "Daily usage retrieved successfully",
  "data": [
    {
      "client_id": "uuid",
      "client_name": "John Doe",
      "date": "2025-01-15",
      "count": 150
    }
  ]
}
```

#### Get Top 3 Clients (Last 24 Hours)
```http
GET /api/usage/top
```

**Response:**
```json
{
  "success": true,
  "message": "Top clients retrieved successfully",
  "data": [
    {
      "client_id": "uuid",
      "client_name": "John Doe",
      "email": "john@example.com",
      "total_requests": 500
    }
  ]
}
```

#### Get Usage Statistics
```http
GET /api/usage/stats
```

#### Get Client Usage
```http
GET /api/usage/client/:client_id
```

#### Get Profile
```http
GET /api/profile
```

#### Refresh Token
```http
POST /api/refresh
```

#### Logout
```http
POST /api/logout
```

## 🔧 Configuration

### Environment Variables

```bash
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=api_tracking
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Rate Limiting
RATE_LIMIT_PER_HOUR=1000

# Cache
CACHE_TTL=3600
```

## 🏗️ Architecture

### Project Structure
```
nexmedis-golang/
├── db/                 # Database connections
│   ├── database.go     # PostgreSQL setup
│   └── redis.go        # Redis setup
├── handler/            # HTTP handlers
│   ├── auth_handler.go
│   ├── client_handler.go
│   ├── log_handler.go
│   └── usage_handler.go
├── model/              # Data models
│   ├── client.go
│   ├── log.go
│   └── request.go
├── router/             # Routes and middleware
│   ├── middleware.go
│   └── router.go
├── store/              # Data access layer
│   ├── client_store.go
│   └── log_store.go
├── utils/              # Utilities
│   ├── crypto.go
│   ├── jwt.go
│   ├── rate_limiter.go
│   ├── response.go
│   └── validator.go
├── main.go             # Entry point
└── README.md
```

### Technology Stack

- **Framework**: Echo v4 (High-performance HTTP framework)
- **ORM**: GORM (Database abstraction)
- **Database**: PostgreSQL (Primary data store)
- **Cache**: Redis (Caching & Pub/Sub)
- **Authentication**: JWT (golang-jwt/jwt/v5)

## 🔐 Security Features

1. **JWT Authentication** - Secure token-based auth for protected endpoints
2. **API Key Validation** - Cryptographic API key generation and validation
3. **Rate Limiting** - Per-client hourly rate limits
4. **Input Validation** - Comprehensive request validation
5. **SQL Injection Protection** - Parameterized queries via GORM
6. **Security Headers** - CORS, XSS, Content-Type protection

## ⚡ Performance Optimizations

1. **Redis Caching**
   - 1-hour TTL for usage endpoints
   - Cache invalidation on data updates
   - Graceful fallback to database

2. **Database Optimizations**
   - Composite indexes on frequently queried fields
   - Batch insert operations for logs
   - Connection pooling (max 100 connections)

## 🧪 Testing

Not yet implemented

## 📦 Deployment

### Production Considerations

1. **Environment Variables**: Use secrets management
2. **Database**: Enable SSL, use read replicas
3. **Redis**: Enable persistence, use Redis Cluster
4. **Monitoring**: Add Prometheus metrics
5. **Logging**: Structured logging with log aggregation
6. **Load Balancing**: Use Nginx or cloud load balancer
7. **HTTPS**: Enable TLS/SSL certificates

## 📊 Database Schema

### Clients Table
```sql
CREATE TABLE clients (
    id UUID PRIMARY KEY,
    client_id VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    api_key VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### API Logs Table
```sql
CREATE TABLE api_logs (
    id UUID PRIMARY KEY,
    client_id UUID REFERENCES clients(id),
    api_key VARCHAR NOT NULL,
    ip VARCHAR NOT NULL,
    endpoint VARCHAR NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP
);

CREATE INDEX idx_api_logs_client_timestamp ON api_logs(client_id, timestamp DESC);
CREATE INDEX idx_api_logs_timestamp ON api_logs(timestamp DESC);
CREATE INDEX idx_api_logs_endpoint ON api_logs(endpoint);
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 🔗 Resources

- [Echo Framework Documentation](https://echo.labstack.com/)
- [GORM Documentation](https://gorm.io/)
- [Go Redis Documentation](https://redis.uptrace.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 📞 Support

For issues and questions, please open an issue on GitHub.

---

Built with ❤️ by Huda
