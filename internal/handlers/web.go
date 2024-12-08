package handlers

import (
	"text/template"

	"github.com/gofiber/fiber/v2"
)

func RenderIndex(c *fiber.Ctx) error {
	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/index.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), nil)
}

func RenderAdd(c *fiber.Ctx) error {
	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/add_log.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), nil)
}
