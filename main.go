package main

import (
	"log"
	"os"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/handlers"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/middleware"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, reading from environment variables")
	}
	// Validate JWT secret early
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
}

func main() {
	app := fiber.New()
	db.Connect()

	// Auto migrate User and Post models
	db.DB.AutoMigrate(&models.User{}, &models.Post{})

	// Public routes
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)
	app.Post("/users/:id/avatar", handlers.UploadAvatar)
	app.Post("/block", handlers.BlockUser)
	app.Post("/logout", handlers.Logout)
	app.Post("/request-reset", handlers.RequestPasswordReset)
	app.Post("/reset-password", handlers.ResetPassword)

	// Protected routes - posts group with JWT middleware
	post := app.Group("/posts", middleware.RequireAuth)

	post.Post("/", handlers.CreatePost)
	post.Get("/", handlers.GetPosts)
	post.Get("/:id", handlers.GetPost)
	post.Put("/:id", handlers.UpdatePost)
	post.Patch("/:id", handlers.PatchPost)
	post.Delete("/:id", handlers.DeletePost)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
