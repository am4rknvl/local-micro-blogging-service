package handlers

import (
	"time"

	"github.com/google/uuid"
	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
)

// VoteInput holds the vote value from client
type VoteInput struct {
	Value int8 `json:"value"` // must be 1 (upvote), -1 (downvote), or 0 (remove vote)
}

// VotePost lets a user cast or change their vote on a post
func VotePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	postID := c.Params("id")

	var input VoteInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.Value != 1 && input.Value != -1 && input.Value != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid vote value"})
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	pid, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	var vote models.Vote
	err = db.DB.Where("user_id = ? AND post_id = ?", uid, pid).First(&vote).Error

	if err != nil {
		// No existing vote, create one if value != 0
		if input.Value == 0 {
			return c.JSON(fiber.Map{"message": "No vote to remove"})
		}
		vote = models.Vote{
			ID:        uuid.New(),
			UserID:    uid,
			PostID:    pid,
			Value:     input.Value,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.DB.Create(&vote).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create vote"})
		}
		return c.JSON(fiber.Map{"message": "Vote cast"})
	}

	// Existing vote found
	if input.Value == 0 {
		// Remove vote
		if err := db.DB.Delete(&vote).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove vote"})
		}
		return c.JSON(fiber.Map{"message": "Vote removed"})
	}

	// Update vote value
	vote.Value = input.Value
	vote.UpdatedAt = time.Now()
	if err := db.DB.Save(&vote).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update vote"})
	}

	return c.JSON(fiber.Map{"message": "Vote updated"})
}

// GetVoteScore returns the net vote score (upvotes - downvotes) for a post
func GetVoteScore(c *fiber.Ctx) error {
	postID := c.Params("id")
	pid, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	var score int64
	err = db.DB.Model(&models.Vote{}).
		Select("COALESCE(SUM(value), 0)").
		Where("post_id = ?", pid).
		Scan(&score).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get vote score"})
	}

	return c.JSON(fiber.Map{"score": score})
}
