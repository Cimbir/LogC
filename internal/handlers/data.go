package handlers

import (
	"LogC/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func GetData(c *fiber.Ctx, _appdata *utils.AppData) error {
	return c.SendString("Hello, World!")
}

func SaveData(c *fiber.Ctx, _appdata *utils.AppData) error {
	return c.SendString("Hello, World!")
}
