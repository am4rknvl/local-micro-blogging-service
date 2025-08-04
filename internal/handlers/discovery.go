package handlers

import (
	"strings"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/gofiber/fiber/v2"
)

func SearchUsers(c *fiber.Ctx) error {
	query := strings.ToLower(c.Query("q"))
	query = "%" + query + "%"

	rows, err := db.DB.Raw(`SELECT id, username, bio FROM users WHERE LOWER(username) LIKE ? LIMIT 20`, query).Rows()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Search failed")
	}
	defer rows.Close()

	var results []fiber.Map
	for rows.Next() {
		var id, username, bio string
		rows.Scan(&id, &username, &bio)
		results = append(results, fiber.Map{
			"id":       id,
			"username": username,
			"bio":      bio,
		})
	}
	return c.JSON(results)
}

func TrendingPosts(c *fiber.Ctx) error {
	rows, err := db.DB.Raw(`SELECT id, content FROM trending_posts LIMIT 10`).Rows()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Trending query failed")
	}
	defer rows.Close()

	var posts []fiber.Map
	for rows.Next() {
		var id, content string
		rows.Scan(&id, &content)
		posts = append(posts, fiber.Map{"id": id, "content": content})
	}
	return c.JSON(posts)
}
