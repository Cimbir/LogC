package handlers

import (
	"LogC/internal/models"
	"LogC/internal/utils"
	"encoding/base64"
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
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	isAdmin := sesh.Get("isAdmin")
	if isAdmin == nil || !isAdmin.(bool) {
		return c.Status(403).SendString("Forbidden")
	}

	// Input JSON template
	var requestData struct {
		Data string `json:"data"`
	}

	// Parse JSON
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	// Decode base64 data
	decodedData, err := base64.StdEncoding.DecodeString(requestData.Data)
	if err != nil {
		return c.Status(400).SendString("Invalid base64 data")
	}
	data := models.LogData{Data: decodedData}

	// Save data
	id, err := _appdata.LogDataCol.Add(data)
	if err != nil {
		return c.Status(500).SendString("Failed to save data")
	}

	// Return id
	return c.JSON(fiber.Map{"id": id})
}
