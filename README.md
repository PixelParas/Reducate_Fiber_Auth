# Fiber Auth Server

This is a Go Fiber backend application providing authentication (Signup, Login) and user management, secured by JWT (JSON Web Tokens) and PostgreSQL.

## Prerequisites

- Go (1.20+)
- PostgreSQL running locally (or adjust the `DB_DSN` in the `.env` file)

## Getting Started

1. Make sure your `.env` file is properly configured.
   ```
   DB_DSN="host=localhost user=postgres password=postgres dbname=fiber_auth port=5432 sslmode=disable"
   JWT_SECRET=your_jwt_secret
   PORT=3000
   ```

2. Start the server (this will run on port 3000 by default):
   ```bash
   go run main.go
   ```

## Example cURL Requests

Here are examples of how to interact with the API endpoints.

### 1. Signup

Register a new user.

```bash
curl -s -X POST http://127.0.0.1:3000/api/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser@example.com", "password":"password123"}'
```

### 2. Login

Login to receive a JWT token.

```bash
curl -s -X POST http://127.0.0.1:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testuser@example.com", "password":"password123"}'
```

**Note**: In the following authenticated requests, replace `YOUR_JWT_TOKEN` with the actual `token` value returned from the Login response.

### 3. Get User Profile (Protected)

Requires a valid JWT.

```bash
curl -s -X GET http://127.0.0.1:3000/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get All Users (Admin Only)

Requires a valid JWT belonging to a user with the `admin` role. 

```bash
curl -s -X GET http://127.0.0.1:3000/api/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
