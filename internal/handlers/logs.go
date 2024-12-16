package handlers

import (
	"LogC/internal/models"
	"LogC/internal/utils"
	"encoding/base64"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetLog(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Get the ID from parameters
	id := c.Params("id")
	if id == "" {
		// Get all logs
		logs, err := _appdata.Logs.GetAll()

		// Reverse the logs
		for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
			logs[i], logs[j] = logs[j], logs[i]
		}
		if err != nil {
			return c.Status(500).SendString("Error getting logs")
		}

		// Return the logs
		return c.JSON(logs)
	} else {
		// Convert the ID to an integer
		logId, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).SendString("Invalid ID")
		}

		// Get the log from the database
		log, err := _appdata.Logs.GetByID(logId)
		if err != nil {
			return c.Status(404).SendString("Log not found")
		}

		// Get the log items from the database
		log_items, err := _appdata.LogItems.GetByField("log_id", id)
		if err != nil {
			return c.Status(500).SendString("Error getting log items")
		}
		log.Items = log_items

		// Return the log
		return c.JSON(log)
	}
}

func SaveLog(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	isAdmin := sesh.Get("isAdmin")
	if isAdmin == nil || !isAdmin.(bool) {
		return c.Status(403).SendString("Forbidden")
	}

	// Parse JSON
	var log models.Log
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	// Save log
	id, err := _appdata.Logs.Add(log)
	if err != nil {
		return c.Status(500).SendString("Failed to save log")
	}

	// Save log items
	for i, item := range log.Items {
		item.LogId = id
		item.Order = i

		// Save Image data
		if item.Type == models.Image {
			// Decode base64 data
			decodedData, err := base64.StdEncoding.DecodeString(item.Content)
			if err != nil {
				return c.Status(400).SendString("Invalid base64 data")
			}
			data := models.LogData{Data: decodedData}

			// Save data
			id, err := _appdata.LogDataCol.Add(data)
			if err != nil {
				return c.Status(500).SendString("Failed to save data")
			}

			item.Content = strconv.Itoa(id)
		}

		_, err := _appdata.LogItems.Add(item)
		if err != nil {
			return c.Status(500).SendString("Failed to save log items")
		}
	}

	// Return id
	return c.JSON(fiber.Map{"id": id})
}

func DeleteLog(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	isAdmin := sesh.Get("isAdmin")
	if isAdmin == nil || !isAdmin.(bool) {
		return c.Status(403).SendString("Forbidden")
	}

	// Get the ID from parameters
	id := c.Params("id")

	// Convert the ID to an integer
	logId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	// Get items from the database
	items, err := _appdata.LogItems.GetByField("log_id", id)
	if err != nil {
		return c.Status(500).SendString("Error getting log items")
	}

	// Delete items
	for _, item := range items {
		// Delete Image data
		if item.Type == models.Image {
			// Convert the ID to an integer
			dataId, err := strconv.Atoi(item.Content)
			if err != nil {
				return c.Status(400).SendString("Invalid ID")
			}

			// Remove data from the database
			err = _appdata.LogDataCol.Remove(dataId)
			if err != nil {
				return c.Status(404).SendString("Data not found")
			}
		}

		// Remove item from the database
		err = _appdata.LogItems.Remove(item.Id)
		if err != nil {
			return c.Status(404).SendString("Item not found")
		}
	}

	// Remove log from the database
	err = _appdata.Logs.Remove(logId)
	if err != nil {
		return c.Status(404).SendString("Log not found")
	}

	// Return the log
	return c.SendStatus(200)
}
