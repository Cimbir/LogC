package handlers

import (
	apiM "LogC/internal/models/api"
	storeM "LogC/internal/models/store"
	"LogC/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Handlers

func GetComments(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Get ID from parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Must have id")
	}

	// Get comments from database
	comments, err := _appdata.Comments.GetByField("log_id", id)
	if err != nil {
		return c.Status(404).SendString("Data not found")
	}

	if comments == nil {
		comments = []storeM.Comment{}
	}

	// Reverse comments
	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
	}

	// Convert to response
	var response []apiM.CommentResponse
	for _, comment := range comments {
		response = append(response, apiM.ToCommentResponse(comment, ""))
	}

	// Return comments
	return c.JSON(response)
}

func SaveComment(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	userId := sesh.Get("userId")
	if userId == nil {
		return c.Status(403).SendString("Forbidden")
	}

	var comment apiM.CommentRequest
	// Parse JSON
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	// Convert to model
	commentModel := apiM.FromCommentRequest(comment)

	// Save comment
	id, err := _appdata.Comments.Add(commentModel)
	if err != nil {
		return c.Status(500).SendString("Failed to save comment")
	}

	// Return id
	return c.JSON(fiber.Map{"id": id})
}
