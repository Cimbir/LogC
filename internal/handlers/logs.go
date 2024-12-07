package handlers

import (
	"LogC/internal/models"
	"LogC/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetLog(c *fiber.Ctx, _appdata *utils.AppData) error {
	if c.Method() != "GET" {
		return c.SendStatus(405)
	}

	id := c.Params("id")
	if id == "" {
		// Return all logs
		logs, err := _appdata.Logs.GetAll()
		if err != nil {
			return c.Status(500).SendString("Error getting logs")
		}
		return c.JSON(logs)
	} else {
		// Return the log with the specified ID
		logId, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).SendString("Invalid ID")
		}
		log, err := _appdata.Logs.GetByID(logId)
		if err != nil {
			return c.Status(404).SendString("Log not found")
		}

		log_items, err := _appdata.LogItems.GetByField("log_id", id)
		if err != nil {
			return c.Status(500).SendString("Error getting log items")
		}
		log.Items = log_items

		return c.JSON(log)
	}
}

func SaveLog(c *fiber.Ctx, _appdata *utils.AppData) error {
	if c.Method() != "POST" {
		return c.SendStatus(405)
	}

	var log models.Log
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	id, err := _appdata.Logs.Add(log)
	if err != nil {
		return c.Status(500).SendString("Failed to save log")
	}

	for i, item := range log.Items {
		item.LogId = id
		item.Order = i
		_, err := _appdata.LogItems.Add(item)
		if err != nil {
			return c.Status(500).SendString("Failed to save log items")
		}
	}

	return c.JSON(fiber.Map{"id": id})
}
