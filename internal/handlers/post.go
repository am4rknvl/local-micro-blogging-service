package handlers

import (
	"time"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Create a new post
func CreatePost(c *fiber.Ctx) error {
	type Input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Simulate user ID from JWT (replace this with c.Locals in real app)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	post := models.Post{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.Status(201).JSON(post)
}

// Get all posts
func GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	if err := db.DB.Order("created_at desc").Find(&posts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch posts"})
	}
	return c.JSON(posts)
}

// Get post by ID
func GetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post
	if err := db.DB.First(&post, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}
	return c.JSON(post)
}

// Update post (PUT = full replace)
func UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")

	var post models.Post
	if err := db.DB.First(&post, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}

	var input models.Post
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	post.Title = input.Title
	post.Content = input.Content
	post.UpdatedAt = time.Now()

	if err := db.DB.Save(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update post"})
	}

	return c.JSON(post)
}

// Partial update (PATCH)
func PatchPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	if err := db.DB.First(&post, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	updates["updated_at"] = time.Now()
	if err := db.DB.Model(&post).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to patch post"})
	}

	return c.JSON(post)
}

// Delete post
func DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&models.Post{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete post"})
	}
	return c.JSON(fiber.Map{"message": "Post deleted"})
}
