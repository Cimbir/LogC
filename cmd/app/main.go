package main

import (
	"LogC/internal/handlers"
	models "LogC/internal/models/store"
	"LogC/internal/store"
	"LogC/internal/utils"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type HandlerFunc func(*fiber.Ctx, *utils.AppData) error

func addAdmin(_appdata *utils.AppData) error {
	// Check if admin exists
	admins, err := _appdata.Users.GetByField("is_admin", true)
	if err != nil {
		return err
	}

	// If no admin exists, create one
	if len(admins) == 0 {
		admin := models.User{
			Username: "Cim",
			Password: "$2a$10$UI5NYV/TPz0MwzKo6AJ2tugz77p3KwUjIP9FPVbmRM.AgeasHOAUO",
			IsAdmin:  true,
		}
		_, err = _appdata.Users.Add(admin)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateEncryptionKey() (string, error) {
	key := make([]byte, 32) // 256 bits
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func main() {
	// Initialize the database
	db, err := store.InitDB(store.DBFilename, store.LogsTable, store.ItemsTable, store.DataTable, store.UserTable, store.CommentTable)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		return
	}
	defer db.Close()

	// Initialize the app data
	var _appdata utils.AppData
	_appdata.Logs = store.NewSQLDB[models.Log](db, store.LogsTable)
	_appdata.LogItems = store.NewSQLDB[models.LogItem](db, store.ItemsTable)
	_appdata.LogDataCol = store.NewSQLDB[models.LogData](db, store.DataTable)
	_appdata.Users = store.NewSQLDB[models.User](db, store.UserTable)
	_appdata.Comments = store.NewSQLDB[models.Comment](db, store.CommentTable)

	// Add admin user
	if err := addAdmin(&_appdata); err != nil {
		fmt.Println("Error adding admin user:", err)
		return
	}

	// Wrapping the handler functions
	handler_wrapper := func(handler HandlerFunc) fiber.Handler {
		return func(c *fiber.Ctx) error {
			return handler(c, &_appdata)
		}
	}

	// Create a new Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Session
	encryptionKey, err := generateEncryptionKey()
	if err != nil {
		fmt.Println("Error generating encryption key:", err)
		return
	}
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptionKey,
	}))
	store := session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Strict",
	})
	app.Use(func(c *fiber.Ctx) error {
		sesh, err := store.Get(c)
		if err != nil {
			return err
		}
		c.Locals("session", sesh)
		return c.Next()
	})

	// Define the routes
	app.Static("/static", "./web/static")

	// Render web
	app.Get("/", handlers.RenderIndex)
	app.Get("/add", handlers.RenderAdd)
	app.Get("/login", handlers.RenderLogin)
	app.Get("/user-management", handlers.RenderUserManagement)
	app.Get("/timeline", handlers.RenderTimeline)
	app.Get("/view/:id", handlers.RenderLogView)
	// Logs
	app.Get("/api/logs/get/:id?", handler_wrapper(handlers.GetLog))
	app.Post("/api/logs/add", handler_wrapper(handlers.SaveLog))
	app.Delete("/api/logs/delete/:id", handler_wrapper(handlers.DeleteLog))
	// Data
	app.Get("/api/data/get/:id", handler_wrapper(handlers.GetData))
	app.Post("/api/data/add", handler_wrapper(handlers.SaveData))
	// Users
	app.Post("/api/users/register", handler_wrapper(handlers.RegisterUser))
	app.Post("/api/users/login", handler_wrapper(handlers.LoginUser))
	app.Get("/api/users/isAdmin", handler_wrapper(handlers.IsAdmin))
	app.Get("/api/users/isLoggedIn", handler_wrapper(handlers.IsLoggedIn))
	app.Get("/api/users/get/:id?", handler_wrapper(handlers.GetUser))
	app.Delete("/api/users/delete/:id", handler_wrapper(handlers.DeleteUser))
	// Comments
	app.Get("/api/comments/get/:id", handler_wrapper(handlers.GetComments))
	app.Post("/api/comments/add", handler_wrapper(handlers.SaveComment))

	// Run the app
	fmt.Println("Server is running at 8090 port.")
	app.Listen(":8090")
}
