package handlers

import (
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	isAdmin := sesh.Get("isAdmin")
	if isAdmin == nil || !isAdmin.(bool) {
		return c.Status(403).SendString("Forbidden")
	}

	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/add_log.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), nil)
}

func RenderLogin(c *fiber.Ctx) error {
	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/login.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), nil)
}

func RenderUserManagement(c *fiber.Ctx) error {
	// Check if user is admin
	sesh := c.Locals("session").(*session.Session)
	isAdmin := sesh.Get("isAdmin")
	if isAdmin == nil || !isAdmin.(bool) {
		return c.Status(403).SendString("Forbidden")
	}

	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/user_management.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), nil)
}

func RenderTimeline(c *fiber.Ctx) error {
	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/timeline.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), nil)
}

func RenderLogView(c *fiber.Ctx) error {
	id := c.Params("id")

	tpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/view.html",
	))
	c.Response().Header.Set("Content-Type", "text/html")
	return tpl.Execute(c.Response().BodyWriter(), fiber.Map{
		"ID": id,
	})
}
