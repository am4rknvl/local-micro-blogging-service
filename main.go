// @title Social Network API
// @version 1.0
// @description This is a microblogging and social networking API.
// @host localhost:8080
// @BasePath /

package main

import (
	"log"
	"os"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/handlers"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/middleware"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/ws"
	"github.com/am4rknvl/local-micro-blogging-service.git/jobs"

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

	// Auto migrate all models
	db.DB.AutoMigrate(&models.User{}, &models.Post{})
	db.DB.AutoMigrate(&models.Follow{})
	db.DB.AutoMigrate(&models.Vote{})
	db.DB.AutoMigrate(&models.Comment{})
	db.DB.AutoMigrate(&models.Message{})
	db.DB.AutoMigrate(&models.Friend{})
	db.DB.AutoMigrate(&models.Notification{})


	// Get the underlying sql.DB for background jobs
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}
	jobs.StartAutoDeleteJob(sqlDB)

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

	// Protected routes - profile group with JWT middleware
	profile := app.Group("/profile", middleware.RequireAuth)
	profile.Get("/:id", handlers.GetProfile)
	profile.Put("/:id", handlers.UpdateProfile)

	// Protected routes - follow group with JWT middleware
	follow := app.Group("/follow", middleware.RequireAuth)
	follow.Post("/:id", handlers.FollowUser)
	follow.Delete("/:id", handlers.UnfollowUser)
	follow.Get("/followers/:id", handlers.GetFollowers)
	follow.Get("/following/:id", handlers.GetFollowing)

	// Protected routes - vote group with JWT middleware
	vote := app.Group("/votes", middleware.RequireAuth)
	vote.Post("/:id", handlers.VotePost)
	vote.Get("/:id/score", handlers.GetVoteScore)

	// Protected routes - comment group with JWT middleware
	comment := app.Group("/posts/:id/comments", middleware.RequireAuth)

	comment.Post("/", handlers.CreateComment)
	comment.Get("/", handlers.GetComments)
	comment.Patch("/:commentId", handlers.UpdateComment)
	comment.Delete("/:commentId", handlers.DeleteComment)

	// Protected routes - messages group with JWT middleware
	messages := app.Group("/conversations/:id/messages", middleware.RequireAuth)
	messages.Post("/", handlers.SendMessage)
	messages.Get("/", handlers.GetMessages)

	app.Post("/admin/delete-old", func(c *fiber.Ctx) error {
		if err := jobs.DeleteOldMessages(sqlDB); err != nil {
			return c.Status(500).SendString("Failed to delete old messages")
		}
		return c.SendString("Old unsaved messages deleted manually")
	})

	app.Patch("/messages/:id/save", handlers.SaveMessage(sqlDB))

	app.Get("/ws/chat/:conversationID", handlers.WebSocketHandler())

	app.Post("/friend-request", handlers.SendFriendRequest)
	app.Post("/respond-request", handlers.RespondToFriendRequest)
	app.Get("/search", handlers.SearchUsers)
	app.Get("/trending", handlers.TrendingPosts)

	app.Get("/friend-requests", handlers.GetFriendRequests)
	app.Get("/friends", handlers.GetFriendTree)
	app.Get("/friend-tree", handlers.GetFriendTree)



	// Start WebSocket manager
	go ws.ManagerInstance.Run()

	jobs.StartMessageCleanupJob()
	// Start server
	log.Fatal(app.Listen(":3000"))
}
