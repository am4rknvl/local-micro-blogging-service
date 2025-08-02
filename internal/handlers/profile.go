package handlers

import (
	"time"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	paramID := c.Params("id")

	var user models.User
	if err := db.DB.First(&user, "id = ?", paramID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Hide profile if private and requester isn’t owner (or admin)
	requesterID, _ := c.Locals("userID").(string)
	if user.IsPrivate && requesterID != paramID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Profile is private"})
	}

	// Clear sensitive data
	user.Password = ""

	return c.JSON(user)
}


func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	paramID := c.Params("id")

	// Only allow users to update their own profile
	if userID != paramID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type ProfileUpdateInput struct {
		Username  *string `json:"username,omitempty"`
		Bio       *string `json:"bio,omitempty"`
		Website   *string `json:"website,omitempty"`
		Location  *string `json:"location,omitempty"`
		IsPrivate *bool   `json:"is_private,omitempty"`
	}

	var input ProfileUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Bio != nil {
		user.Bio = *input.Bio
	}
	if input.Website != nil {
		user.Website = *input.Website
	}
	if input.Location != nil {
		user.Location = *input.Location
	}
	if input.IsPrivate != nil {
		user.IsPrivate = *input.IsPrivate
	}

	user.UpdatedAt = time.Now()

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.JSON(user)
}
