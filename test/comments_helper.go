package test

import (
	"LogC/internal/handlers"
	storeM "LogC/internal/models/store"
	"LogC/internal/utils"
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

type setuphandler func(*testing.T) (utils.AppData, teardown)

func GetCommentsTestHelper(t *testing.T, stp setuphandler) {
	appData, td := stp(t)
	defer td()

	app := fiber.New()

	// Add a log to the database for testing
	log := storeM.Log{Title: "Test Log", Date: time.Now()}
	logId, err := appData.Logs.Add(log)
	if err != nil {
		t.Fatalf("Failed to add log: %v", err)
	}

	// Add user to the database for testing
	user := storeM.User{Username: "Test	User", Password: "password"}
	userId, err := appData.Users.Add(user)
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}

	// Add comment to the database for testing
	comment := storeM.Comment{LogId: logId, UserId: userId, Content: "Test Comment"}
	commentId, err := appData.Comments.Add(comment)
	if err != nil {
		t.Fatalf("Failed to add comment: %v", err)
	}

	// Define the route
	app.Get("/api/comments/get/:id", func(c *fiber.Ctx) error {
		return handlers.GetComments(c, &appData)
	})

	// Create a request to get the log
	req := httptest.NewRequest("GET", "/api/comments/get/"+strconv.Itoa(logId), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response
	var returnedComments []storeM.Comment
	if err := json.NewDecoder(resp.Body).Decode(&returnedComments); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, 1, len(returnedComments))
	assert.Equal(t, commentId, returnedComments[0].Id)
}

func SaveCommentTestHelper(t *testing.T, stp setuphandler) {
	appData, td := stp(t)
	defer td()

	// Add log to the database for testing
	log := storeM.Log{Title: "Test Log", Date: time.Now()}
	logId, err := appData.Logs.Add(log)
	if err != nil {
		t.Fatalf("Failed to add log: %v", err)
	}

	// Add user to the database for testing
	user := storeM.User{Username: "Test User", Password: "password"}
	userId, err := appData.Users.Add(user)
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}

	app := fiber.New()

	// Define the route
	store := session.New()
	app.Post("/api/comments/add", func(c *fiber.Ctx) error {
		sesh, err := store.Get(c)
		if err != nil {
			return err
		}
		sesh.Set("userId", userId)
		c.Locals("session", sesh)
		return handlers.SaveComment(c, &appData)
	})

	// Create a comment to save
	comment := storeM.Comment{LogId: logId, UserId: userId, Content: "Test Comment"}
	commentJSON, err := json.Marshal(comment)
	if err != nil {
		t.Fatalf("Failed to marshal comment: %v", err)
	}

	// Create a request to save comment
	req := httptest.NewRequest("POST", "/api/comments/add", bytes.NewReader(commentJSON))
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

	// Verify the comment was saved
	savedLog, err := appData.Comments.GetByID(id)
	if err != nil {
		t.Fatalf("Failed to get saved comment: %v", err)
	}

	assert.Equal(t, comment.Content, savedLog.Content)
}
