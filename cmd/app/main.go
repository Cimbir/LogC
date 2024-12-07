package main

import (
	"LogC/internal/handlers"
	"LogC/internal/models"
	"LogC/internal/store"
	"LogC/internal/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type HandlerFunc func(*fiber.Ctx, *utils.AppData) error

func main() {
	// Change the working directory to the root of the project
	err := os.Chdir(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error changing working directory:", err)
		return
	}

	// Initialize the database
	db, err := store.InitDB(store.DBFilename, store.LogsTable, store.ItemsTable, store.DataTable)
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

	// Wrapping the handler functions
	handler_wapper := func(handler HandlerFunc) fiber.Handler {
		return func(c *fiber.Ctx) error {
			return handler(c, &_appdata)
		}
	}

	// Create a new Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	//app.Use(recover.New())

	// Define the routes
	app.Static("/static", "./web/static")

	// Logs
	app.Get("/api/logs/get/:id?", handler_wapper(handlers.GetLog))
	app.Post("/api/logs/add", handler_wapper(handlers.SaveLog))

	// Data
	app.Get("/api/data/get/:id", handler_wapper(handlers.GetData))
	app.Post("/api/data/add", handler_wapper(handlers.SaveData))

	// Run the app
	fmt.Println("Server is running at 8090 port.")
	app.Listen(":8090")
}
