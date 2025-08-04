package handlers

import (
	"time"

	"github.com/google/uuid"
	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
)


type EnrichedComment struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsOP      bool      `json:"is_op"`

	// User Info (Embedded)
	User struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Avatar   string    `json:"avatar"`
	} `json:"user"`
}

// CreateComment creates a comment on a post
func CreateComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	postID := c.Params("id")

	var input struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&input); err != nil || input.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	pid, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	comment := models.Comment{
		ID:        uuid.New(),
		PostID:    pid,
		UserID:    uid,
		Content:   input.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	return c.JSON(comment)
}

// GetComments gets all comments for a post
func GetComments(c *fiber.Ctx) error {
	postID := c.Params("id")

	pid, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	// Get post to find OP
	var post models.Post
	if err := db.DB.First(&post, "id = ?", pid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	// Get all comments on that post
	var comments []models.Comment
	if err := db.DB.Where("post_id = ?", pid).Order("created_at asc").Find(&comments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comments"})
	}

	// Prepare enriched response
	var enriched []EnrichedComment
	for _, cmt := range comments {
		var user models.User
		if err := db.DB.Select("id", "username", "avatar").First(&user, "id = ?", cmt.UserID).Error; err != nil {
			continue // silently skip if user not found
		}

		enriched = append(enriched, EnrichedComment{
			ID:        cmt.ID,
			Content:   cmt.Content,
			CreatedAt: cmt.CreatedAt,
			IsOP:      cmt.UserID == post.UserID,
			User: struct {
				ID       uuid.UUID `json:"id"`
				Username string    `json:"username"`
				Avatar   string    `json:"avatar"`
			}{
				ID:       user.ID,
				Username: user.Username,
				Avatar:   user.Avatar,
			},
		})
	}

	return c.JSON(enriched)
}

// UpdateComment updates a comment (only by owner)
func UpdateComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	commentID := c.Params("commentId")

	var input struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&input); err != nil || input.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	cid, err := uuid.Parse(commentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	var comment models.Comment
	if err := db.DB.First(&comment, "id = ?", cid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	if comment.UserID.String() != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	comment.Content = input.Content
	comment.UpdatedAt = time.Now()

	if err := db.DB.Save(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	return c.JSON(comment)
}

// DeleteComment deletes a comment (only by owner/Op)
func DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	commentID := c.Params("commentId")

	cid, err := uuid.Parse(commentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	var comment models.Comment
	if err := db.DB.First(&comment, "id = ?", cid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	if comment.UserID.String() != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := db.DB.Delete(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete comment"})
	}

	return c.JSON(fiber.Map{"message": "Comment deleted"})
}
