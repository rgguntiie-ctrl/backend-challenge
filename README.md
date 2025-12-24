## Architecture

This project follows Clean Architecture principles with clear separation of concerns:
```
backend-challenge/
├── cmd/backend-api/          # Application entry point
├── config/                   # Configuration management
├── infrastructure/           # External services (DB, Redis, JWT)
├── internal/
│   ├── core/
│   │   ├── domain/          # Business entities
│   │   ├── ports/           # Interfaces (contracts)
│   │   └── services/        # Business logic
│   └── adapters/
│       ├── handlers/        # HTTP handlers
│       ├── repositories/    # Data access layer
│       └── cache/          # Cache implementations
├── middlewares/             # HTTP middlewares
└── docs/                   # Swagger documentation
```

## 🛠️ Technology Stack

- **Language:** Go 1.24
- **Web Framework:** Fiber v2
- **Database:** PostgreSQL 16
- **Cache:** Redis 7
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt/jwt)
- **Documentation:** Swagger (swaggo)
- **Logging:** Zap, Logrus

## 📋 Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose
- Make (optional, for convenience)

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd backend-challenge
```

### 2. Start services with Docker Compose

```bash
docker compose up -d
```

This will start:
- PostgreSQL on port `5432`
- Redis on port `6379`


### 3. Set up environment variables

Copy the example file:
```bash
cp env.example .env
```

The `.env` file should contain:
```env
JWT_SECRET=test-backend-challenge-secret

PSQL_HOST=localhost
PSQL_PORT=5432
PSQL_USER=postgres
PSQL_PASS=123456
PSQL_DB=be_db

REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 4. Install dependencies

```bash
go mod download
```

### 5. Run the application

Using Make:
```bash
make run.backend-api
```

Or directly:
```bash
go run cmd/backend-api/main.go
```

The server will start on `http://localhost:3000`

## 📚 API Documentation

Once the server is running, access the Swagger documentation at:

```
http://localhost:3000/docs/index.html
```

### Available Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | Login and get tokens | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh access token | ❌ |
| POST | `/api/v1/auth/logout` | Logout (revoke tokens) | ✅ |
| GET | `/api/v1/users/me` | Get current user profile | ✅ |
| GET | `/health` | Health check | ❌ |


## 🧪 API Usage Examples

### 1. Register a new user

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "message": "registered"
}
```

### 2. Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### 3. Get user profile (Protected)

```bash
curl -X GET http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com"
}
```

### 4. Refresh access token

```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 5. Logout (Revoke tokens)

```bash
curl -X POST http://localhost:3000/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "message": "logged out successfully"
}
```
