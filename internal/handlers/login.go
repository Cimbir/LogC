package handlers

import (
	"LogC/internal/models"
	"LogC/internal/services"
	"LogC/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Get user json data from the request body, add it to the database and log in with session
// The username must be unique
func RegisterUser(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Parse JSON
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	// Check if username or password is empty
	if user.Username == "" || user.Password == "" {
		return c.Status(400).SendString("Username or password is empty")
	}

	// Check if username already exists
	sameUsernames, err := _appdata.Users.GetByField("username", user.Username)
	if err != nil || len(sameUsernames) > 0 {
		return c.Status(400).SendString("Username already exists")
	}

	// Hash password
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).SendString("Error hashing password")
	}
	user.Password = hashedPassword
	user.IsAdmin = false

	// Add user to database
	id, err := _appdata.Users.Add(user)
	if err != nil {
		return c.Status(500).SendString("Error adding user")
	}

	// Store user info in session
	sesh := c.Locals("session").(*session.Session)
	sesh.Set("userId", id)
	sesh.Set("isAdmin", false)
	if err := sesh.Save(); err != nil {
		return c.Status(500).SendString("Error saving session")
	}

	return c.JSON(fiber.Map{"id": id})
}

// Log into account with session given username and password.
func LoginUser(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Parse JSON
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	// Check if user exists
	storedUser, err := _appdata.Users.GetByField("username", user.Username)
	if err != nil || len(storedUser) == 0 {
		return c.Status(404).SendString("User not found")
	}

	// Check password
	if !services.CheckPasswordHash(user.Password, storedUser[0].Password) {
		return c.Status(401).SendString("Invalid password")
	}

	// Store user info in session
	sesh := c.Locals("session").(*session.Session)
	sesh.Set("userId", storedUser[0].Id)
	sesh.Set("isAdmin", storedUser[0].IsAdmin)
	if err := sesh.Save(); err != nil {
		return c.Status(500).SendString("Error saving session")
	}

	return c.JSON(fiber.Map{"message": "Logged in"})
}
