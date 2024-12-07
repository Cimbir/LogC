package handlers

import (
	"LogC/internal/models"
	"LogC/internal/utils"
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

	return c.JSON(data)
}

func SaveData(c *fiber.Ctx, _appdata *utils.AppData) error {
	if c.Method() != "POST" {
		return c.SendStatus(405)
	}

	var data models.LogData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	id, err := _appdata.LogDataCol.Add(data)
	if err != nil {
		return c.Status(500).SendString("Failed to save data")
	}

	return c.JSON(fiber.Map{"id": id})
}
