package handlers

import (
	models "LogC/internal/models/store"
	"LogC/internal/utils"
	"fmt"
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetData(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Get ID from parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Must have id")
	}

	// Convert id to int
	dataId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// Get data from database
	data, err := _appdata.LogDataCol.GetByID(dataId)
	if err != nil {
		return c.Status(404).SendString("Data not found")
	}

	// Return data
	c.Set("Content-Type", "image/png")
	return c.Send(data.Data)
}

// input - json with field "data" containing base64 encoded data
// saves in byte array format
func SaveData(c *fiber.Ctx, _appdata *utils.AppData) error {
	fmt.Println("ENTERED")

	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	isAdmin := sesh.Get("isAdmin")
	if isAdmin == nil || !isAdmin.(bool) {
		return c.Status(403).SendString("Forbidden")
	}

	fmt.Println("first")
	// Get file from form-data
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Failed to get file from form-data")
	}

	fmt.Println("second")
	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Failed to open file")
	}
	defer fileContent.Close()

	fmt.Println("third")
	// Read file content
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(500).SendString("Failed to read file content")
	}

	fmt.Println("fourth")
	data := models.LogData{Data: fileBytes}

	// Save data
	id, err := _appdata.LogDataCol.Add(data)
	if err != nil {
		return c.Status(500).SendString("Failed to save data")
	}

	// Return id
	return c.JSON(fiber.Map{"id": id})
}
