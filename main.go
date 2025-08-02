package main

import (
	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/handlers"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db.Connect()

	
	db.DB.AutoMigrate(&models.User{})

	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)
	app.Post("/users/:id/avatar", handlers.UploadAvatar)
	app.Post("/block", handlers.BlockUser)
	app.Post("/logout", handlers.Logout)
	app.Post("/request-reset", handlers.RequestPasswordReset)
	app.Post("/reset-password", handlers.ResetPassword)




	app.Listen(":3000")
}
