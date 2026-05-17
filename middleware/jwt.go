package middleware

import (
	"fmt"
	"strings"

	"fiber_auth/config" // Adjust this if your actual module name is different

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var cfg *config.Config

// InitMiddleware injects the configuration into the middleware package
func InitMiddleware(c *config.Config) {
	cfg = c
}

// RequireAuth protects routes by verifying a JWT token
func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// Ensure the header isn't empty and starts with "Bearer "
	// Note: Added the space after "Bearer " to be more precise
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or malformed JWT",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate the JWT Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC (the standard for standard secret keys)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key from our config
		return []byte(cfg.JWTSecret), nil
	})

	// Handle parsing errors (e.g., token is expired, malformed, or signature is invalid)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Extract claims (the payload) from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to parse token claims",
		})
	}

	// Extract the user ID
	// Important: JSON numbers are always decoded as float64 by default in Go.
	// We must assert it as float64 first, then cast it to uint.
	idFloat, ok := claims["id"].(float64) // Make sure "id" matches the key you used when creating the token!
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}
	userID := uint(idFloat)

	// Extract the role
	role, ok := claims["role"].(string) // Make sure "role" matches the key you used when creating the token!
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid role claim in token",
		})
	}

	// Store the extracted, strongly-typed values in Fiber's context.
	// This allows subsequent route handlers to easily retrieve the user ID and role.
	c.Locals("user_id", userID)
	c.Locals("role", role)

	// Proceed to the next middleware or route handler
	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)

	if !ok || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access reuired",
		})
	}

	return c.Next()
}
