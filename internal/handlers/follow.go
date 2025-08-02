package handlers

import (
	"time"

	"github.com/google/uuid"
	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
)

// FollowUser lets the authenticated user follow another user
func FollowUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)   // current user
	followeeID := c.Params("id")            // user to follow

	if userID == followeeID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot follow yourself"})
	}

	// Check if already following
	var existing models.Follow
	err := db.DB.Where("follower_id = ? AND followee_id = ?", userID, followeeID).First(&existing).Error
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already following"})
	}

	follow := models.Follow{
		ID:         uuid.New(),
		FollowerID: uuid.MustParse(userID),
		FolloweeID: uuid.MustParse(followeeID),
		CreatedAt:  time.Now(),
	}

	if err := db.DB.Create(&follow).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to follow user"})
	}

	return c.JSON(fiber.Map{"message": "Successfully followed user"})
}

// UnfollowUser lets the authenticated user unfollow another user
func UnfollowUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)   // current user
	followeeID := c.Params("id")            // user to unfollow

	if userID == followeeID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot unfollow yourself"})
	}

	err := db.DB.Where("follower_id = ? AND followee_id = ?", userID, followeeID).Delete(&models.Follow{}).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unfollow user"})
	}

	return c.JSON(fiber.Map{"message": "Successfully unfollowed user"})
}

// GetFollowers returns a list of users who follow the given user
func GetFollowers(c *fiber.Ctx) error {
	userID := c.Params("id")

	var followers []models.User
	err := db.DB.Raw(`
		SELECT u.* FROM users u
		JOIN follows f ON f.follower_id = u.id
		WHERE f.followee_id = ?
	`, userID).Scan(&followers).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get followers"})
	}

	return c.JSON(followers)
}

// GetFollowing returns a list of users the given user is following
func GetFollowing(c *fiber.Ctx) error {
	userID := c.Params("id")

	var following []models.User
	err := db.DB.Raw(`
		SELECT u.* FROM users u
		JOIN follows f ON f.followee_id = u.id
		WHERE f.follower_id = ?
	`, userID).Scan(&following).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get following"})
	}

	return c.JSON(following)
}
