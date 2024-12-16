package handlers

import (
	"LogC/internal/models"
	"LogC/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetComments(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Get ID from parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Must have id")
	}

	// Convert id to int
	logId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// Get comments from database
	comments, err := _appdata.Comments.GetByField("log_id", logId)
	if err != nil {
		return c.Status(404).SendString("Data not found")
	}

	// Return comments
	return c.JSON(comments)
}

// input json of comment
// saves in database
func SaveComment(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	userId := sesh.Get("userId")
	if userId == nil {
		return c.Status(403).SendString("Forbidden")
	}

	var comment models.Comment

	// Parse JSON
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}
	comment.Date = time.Now()
	comment.UserId = userId.(int)

	// Save comment
	id, err := _appdata.Comments.Add(comment)
	if err != nil {
		return c.Status(500).SendString("Failed to save comment")
	}

	// Return id
	return c.JSON(fiber.Map{"id": id})
}
