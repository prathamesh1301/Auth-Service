# 🔐 Auth Service

A lightweight authentication microservice built in Go, featuring JWT-based authentication, bcrypt password hashing, and PostgreSQL storage. Dockerized for easy deployment.

## Architecture

```
cmd/                        → Application entrypoint & HTTP handlers
├── main.go                 → Bootstrap (config, DB, DI)
├── app.go                  → Router setup, server start
├── auth_handler.go         → Login, Signup, JWT middleware
├── refresh_token_handler.go→ Refresh token handler
└── health.go               → Health check (protected)

internals/
├── db/
│   └── db.go               → PostgreSQL connection pool
├── jwt/
│   └── jwt.go              → Token generation & validation
└── store/
    ├── store.go             → Data access aggregator
    ├── user.go              → User repository (SQL queries)
    └── refresh_token.go     → Refresh token repository

migrations/                 → SQL migration files (golang-migrate)
```

### Design Principles

- **Layered architecture**: `Handler → Store → Database` with clean separation of concerns
- **Interface-driven**: `UserRepository` and `TokenService` interfaces enable easy testing and swappable implementations
- **Implicit interface satisfaction**: Go's duck typing — no `implements` keyword needed
- **Context propagation**: All DB queries accept `context.Context` for timeout/cancellation support
- **Environment-based config**: Secrets and connection strings read from environment variables

## Tech Stack

| Component         | Technology                     |
|-------------------|--------------------------------|
| Language          | Go 1.25                        |
| Router            | [chi](https://github.com/go-chi/chi) |
| Database          | PostgreSQL 15                  |
| Auth              | JWT (HS256) + bcrypt           |
| Migrations        | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Containerization  | Docker + Docker Compose        |

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (for migrations)

### Run with Docker Compose

```bash
# Start PostgreSQL + API
docker compose up --build

# In a separate terminal, run migrations
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```

### Run Locally (without Docker)

```bash
# Make sure PostgreSQL is running on localhost:5432

# Set environment variables (optional, defaults provided)
set DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
set JWT_SECRET=my-super-secret-key

# Run migrations
migrate -path ./migrations -database "%DATABASE_URL%" up

# Start the server
cd cmd
go run .
```

The server will start on `http://localhost:8080`.

## Environment Variables

| Variable       | Description                    | Default                                                            |
|----------------|--------------------------------|--------------------------------------------------------------------|
| `DATABASE_URL` | PostgreSQL connection string   | `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable` |
| `JWT_SECRET`   | Secret key for signing JWTs    | `dev-secret-change-me-in-production`                               |

## API Endpoints

### Public Routes

#### `POST /signup`
Create a new user account.

```bash
curl.exe -X POST http://localhost:8080/signup ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"john\", \"password\": \"mypassword\"}"
```

**Response:** `201 Created`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "12345abcdef..."
}
```

---

#### `POST /login`
Authenticate and receive a JWT token.

```bash
curl.exe -X POST http://localhost:8080/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"john\", \"password\": \"mypassword\"}"
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "12345abcdef..."
}
```

---

#### `POST /refreshToken`
Generate a new access token using a valid refresh token.

```bash
curl.exe -X POST http://localhost:8080/refreshToken ^
  -H "Content-Type: application/json" ^
  -d "{\"refresh_token\": \"<your-refresh-token>\"}"
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "<your-refresh-token>"
}
```

### Protected Routes

These require a valid JWT in the `Authorization` header.

#### `GET /health`
Health check endpoint.

```bash
curl.exe -X GET http://localhost:8080/health ^
  -H "Authorization: Bearer <your-jwt-token>"
```

**Response:** `200 OK` — `Good health`

## Database Schema

```sql
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id TEXT REFERENCES users(username),
    token TEXT UNIQUE,
    expires_at TIMESTAMP
);
```

