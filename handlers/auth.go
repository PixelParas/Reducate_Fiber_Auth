package handlers

import (
	"fiber_auth/config"
	"fiber_auth/database"
	"fiber_auth/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var cfg *config.Config

func InitHandlers(c *config.Config) {
	cfg = c
}

// Signup POST /api/Signup
// 1. Parse the Json into AuthInput
// 2. Validate AuthInput
// 3. Hash using bcrypt
// 4. Save the user to database
// 5. Return the created user
func Signup(c *fiber.Ctx) error {
	var input models.AuthInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := input.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     models.RoleUser,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "An account with this email already exists",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created succesfully",
		"user":    user.ToResponse(),
	})
}

// Login api/login
// 1. Parse and validate the json
// 2. Look up the email
// 3. compare the password with passworhHash
// 4. if valid return a JWT
func Login(c *fiber.Ctx) error {
	var input models.AuthInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := input.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	token, err := generateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT Token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}

func generateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.JWTSecret))
}
