package handlers

import (
	"fmt"
	"time"

	"github.com/am4rknvl/local-micro-blogging-service.git/internal/config"
	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	type SignupInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input SignupInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user := models.User{
		ID:        uuid.New(),
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.JSON(fiber.Map{"message": "Signup successful"})
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if user.Blocked {
		return c.Status(403).JSON(fiber.Map{"error": "User is blocked"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 days expiration, reduce it by anychance...
	})

	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not login"})
	}

	// Return token to client
	return c.JSON(fiber.Map{"token": tokenString})
}


func Logout(c *fiber.Ctx) error {
	// Clear token cookie or just tell client to delete token
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}


func BlockUser(c *fiber.Ctx) error {
	type BlockInput struct {
		UserID string `json:"user_id"`
		Block  bool   `json:"block"`
	}

	var input BlockInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", input.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	user.Blocked = input.Block

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update user"})
	}

	return c.JSON(fiber.Map{"message": "User block status updated"})
}


var resetTokens = make(map[string]string) // map[email]token


func RequestPasswordReset(c *fiber.Ctx) error {
	type ResetRequest struct {
		Email string `json:"email"`
	}

	var req ResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// don't reveal user existence
		return c.JSON(fiber.Map{"message": "If email exists, reset token sent"})
	}

	token := uuid.New().String()
	resetTokens[req.Email] = token

	// In real app: send email with token here
	// For MVP: just return token in response
	return c.JSON(fiber.Map{"reset_token": token})
}


func ResetPassword(c *fiber.Ctx) error {
	type ResetInput struct {
		Email       string `json:"email"`
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
	}

	var input ResetInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	storedToken, ok := resetTokens[input.Email]
	if !ok || storedToken != input.ResetToken {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid reset token"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user.Password = string(hashed)
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update password"})
	}

	// remove token so it can’t be reused
	delete(resetTokens, input.Email)

	return c.JSON(fiber.Map{"message": "Password reset successful"})
}


func UploadAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	savePath := fmt.Sprintf("./uploads/%s_%s", userID, file.Filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not save file"})
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	user.Avatar = savePath
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update avatar"})
	}

	return c.JSON(fiber.Map{"message": "Avatar uploaded", "path": savePath})
}
