package handlers

import (
	apiM "LogC/internal/models/api"
	storeM "LogC/internal/models/store"
	"LogC/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Handlers

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

		// Convert to response
		response := []apiM.LogResponse{}
		for _, log := range logs {
			response = append(response, apiM.ToLogResponse(log))
		}

		// Return the logs
		return c.JSON(response)
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
		logItems, err := _appdata.LogItems.GetByField("log_id", id)
		if err != nil {
			return c.Status(500).SendString("Error getting log items")
		}

		// Get the comments from the database
		comments, err := _appdata.Comments.GetByField("log_id", id)
		if err != nil {
			return c.Status(500).SendString("Error getting comments")
		}

		// Get the usernames from the database
		usernames := map[int]string{}
		for _, comment := range comments {
			if _, ok := usernames[comment.UserId]; ok {
				continue
			}
			user, err := _appdata.Users.GetByID(comment.UserId)
			if err != nil {
				return c.Status(500).SendString("Error getting user")
			}
			usernames[comment.UserId] = user.Username
		}

		// Reverse the comments
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}

		// Create a new response
		response := apiM.ToFullLogResponse(log, logItems, comments, usernames)

		// Return the log
		return c.JSON(response)
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
	var log apiM.LogRequest
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).SendString("Cannot parse JSON")
	}

	// Convert to model
	logModel := apiM.FromLogRequest(log)

	// Save log
	id, err := _appdata.Logs.Add(logModel)
	if err != nil {
		return c.Status(500).SendString("Failed to save log")
	}

	// Save log items
	for i, item := range log.Items {
		itemModel := apiM.FromLogItemRequest(item, i, id)

		_, err := _appdata.LogItems.Add(itemModel)
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
		if item.Type == storeM.Image {
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

func GetCategoryLogs(c *fiber.Ctx, _appdata *utils.AppData) error {
	// Get the amount to get
	amount := c.Params("amount")
	if amount == "" {
		amount = "10"
	}

	// Convert the amount to an integer
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return c.Status(400).SendString("Invalid amount")
	}

	// Get all logs
	logs, err := _appdata.Logs.GetAll()
	if err != nil {
		return c.Status(500).SendString("Error getting logs")
	}

	// Reverse the logs
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	// Convert to response
	response := apiM.ToLogByCategoryResponse(logs, amountInt)

	// Return the logs
	return c.JSON(response)
}
