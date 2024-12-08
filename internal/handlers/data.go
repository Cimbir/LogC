package handlers

import (
	"LogC/internal/models"
	"LogC/internal/utils"
	"encoding/base64"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetData(c *fiber.Ctx, _appdata *utils.AppData) error {
	if c.Method() != "GET" {
		return c.SendStatus(405)
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Must have id")
	}

	dataId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	data, err := _appdata.LogDataCol.GetByID(dataId)
	if err != nil {
		return c.Status(404).SendString("Data not found")
	}
	c.Set("Content-Type", "image/png")
	return c.Send(data.Data)
}

// input - json with field "data" containing base64 encoded data
// saves in byte array format
func SaveData(c *fiber.Ctx, _appdata *utils.AppData) error {
	if c.Method() != "POST" {
		return c.SendStatus(405)
	}

	var requestData struct {
		Data string `json:"data"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	decodedData, err := base64.StdEncoding.DecodeString(requestData.Data)
	if err != nil {
		return c.Status(400).SendString("Invalid base64 data")
	}
	data := models.LogData{Data: decodedData}

	id, err := _appdata.LogDataCol.Add(data)
	if err != nil {
		return c.Status(500).SendString("Failed to save data")
	}

	return c.JSON(fiber.Map{"id": id})
}
