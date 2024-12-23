package test

import (
	"LogC/internal/handlers"
	apiM "LogC/internal/models/api"
	storeM "LogC/internal/models/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/stretchr/testify/assert"
)

func GetLogTestHelper(t *testing.T, stp setuphandler) {
	appData, td := stp(t)
	defer td()

	app := fiber.New()

	// Add a log to the database for testing
	log := storeM.Log{Title: "Test Log", Date: time.Now()}
	id, err := appData.Logs.Add(log)
	if err != nil {
		t.Fatalf("Failed to add log: %v", err)
	}

	// Add items to the database for testing
	items := []storeM.LogItem{
		{LogId: id, Type: storeM.Title, Content: "Description 1"},
		{LogId: id, Type: storeM.Text, Content: "Description 2"},
	}
	for i, item := range items {
		_, err := appData.LogItems.Add(item)
		if err != nil {
			t.Fatalf("Failed to add item %d: %v", i, err)
		}
	}

	// Define the route
	app.Get("/api/logs/get/:id?", func(c *fiber.Ctx) error {
		return handlers.GetLog(c, &appData)
	})

	// Create a request to get the log
	req := httptest.NewRequest("GET", "/api/logs/get/"+strconv.Itoa(id), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var returnedLog apiM.LogResponse
	if err := json.NewDecoder(resp.Body).Decode(&returnedLog); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, log.Title, returnedLog.Title)
	assert.Equal(t, log.Date.Format("2006-01-01"), returnedLog.Date.Format("2006-01-01"))
	assert.Equal(t, len(items), len(returnedLog.Items))

	// Test getting all logs

	// Add new log
	log = storeM.Log{Title: "Test Log 2", Date: time.Now()}
	id, err = appData.Logs.Add(log)
	if err != nil {
		t.Fatalf("Failed to add log: %v", err)
	}

	// Create a request to get the log
	req = httptest.NewRequest("GET", "/api/logs/get/", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var returnedLogs []apiM.LogResponse
	if err := json.NewDecoder(resp.Body).Decode(&returnedLogs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, len(returnedLogs), 2)
}

func SaveLogTestHelper(t *testing.T, stp setuphandler) {
	appData, td := stp(t)
	defer td()

	app := fiber.New()

	// Define the route
	store := session.New()
	app.Post("/api/logs/save", func(c *fiber.Ctx) error {
		sesh, err := store.Get(c)
		if err != nil {
			return err
		}
		sesh.Set("isAdmin", true)
		c.Locals("session", sesh)
		return handlers.SaveLog(c, &appData)
	})

	// Create a log to save
	log := apiM.LogRequest{Title: "New Log", Items: []apiM.LogItemRequest{
		{Type: "Title", Content: "Title"},
		{Type: "Text", Content: "Description"},
	}}
	logJSON, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal log: %v", err)
	}

	// Create a request to save the log
	req := httptest.NewRequest("POST", "/api/logs/save", bytes.NewReader(logJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseMap map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	id, ok := responseMap["id"]
	if !ok {
		t.Fatalf("Response does not contain id")
	}

	// Verify the log was saved
	savedLog, err := appData.Logs.GetByID(id)
	if err != nil {
		t.Fatalf("Failed to get saved log: %v", err)
	}

	savedItems, err := appData.LogItems.GetByField("log_id", strconv.Itoa(id))
	if err != nil {
		t.Fatalf("Failed to get saved items: %v", err)
	}

	assert.Equal(t, log.Title, savedLog.Title)
	assert.Equal(t, len(log.Items), len(savedItems))
}
