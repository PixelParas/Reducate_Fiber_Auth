package handlers

import (
	"fiber_auth/database"
	"fiber_auth/models"

	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authnticaion context",
		})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No User Found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user.ToResponse(),
	})

}

func GetAllUsers(c *fiber.Ctx) error {

	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}

	return c.JSON(fiber.Map{
		"users": responses,
		"count": len(responses),
	})

}

