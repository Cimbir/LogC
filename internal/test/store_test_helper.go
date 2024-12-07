package test

import (
	"LogC/internal/models"
	"LogC/internal/store"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type teardown func()
type setup func(*testing.T) (store.DB[models.Log], store.DB[models.LogItem], store.DB[models.LogData], teardown)

func AddTestHelper(t *testing.T, stp setup) {
	logs, logItems, logData, td := stp(t)
	defer td()

	// log
	log := models.Log{Title: "Test Log", Date: time.Now()}
	id, err := logs.Add(log)
	if err != nil || id == -1 {
		t.Fatalf("Failed to add log: %v", err)
	}

	// item
	item := models.LogItem{LogId: id, Type: models.Text, Content: "Test Item", Order: 1}
	id, err = logItems.Add(item)
	if err != nil || id == -1 {
		t.Fatalf("Failed to add log item: %v", err)
	}

	// data
	data := models.LogData{Data: []byte("Test Data")}
	id, err = logData.Add(data)
	if err != nil || id == -1 {
		t.Fatalf("Failed to add log data: %v", err)
	}
}

func GetAllTestHelper(t *testing.T, stp setup) {
	logs, logItems, logData, td := stp(t)
	defer td()

	// log
	log := models.Log{Title: "Test Log", Date: time.Now()}
	id, err := logs.Add(log)
	allLogs, err := logs.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all logs: %v", err)
	}
	if len(allLogs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(allLogs))
	}

	// item
	item := models.LogItem{LogId: id, Type: models.Text, Content: "Test Item", Order: 1}
	id, err = logItems.Add(item)
	allItems, err := logItems.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all log items: %v", err)
	}
	if len(allItems) != 1 {
		t.Fatalf("Expected 1 log item, got %d", len(allItems))
	}

	// data
	data := models.LogData{Data: []byte("Test Data")}
	id, err = logData.Add(data)
	allData, err := logData.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all log data: %v", err)
	}
	if len(allData) != 1 {
		t.Fatalf("Expected 1 log data, got %d", len(allData))
	}
}

func GetByIDTestHelper(t *testing.T, stp setup) {
	logs, logItems, logData, td := stp(t)
	defer td()

	// log
	log := models.Log{Title: "Test Log", Date: time.Now()}
	id, err := logs.Add(log)
	got_log, err := logs.GetByID(id)
	if err != nil || got_log.Id != id {
		t.Fatalf("Expected log with ID %d, got %v", id, got_log)
	}

	// item
	item := models.LogItem{LogId: id, Type: models.Text, Content: "Test Item", Order: 1}
	id, err = logItems.Add(item)
	got_item, err := logItems.GetByID(id)
	if err != nil || got_item.Id != id {
		t.Fatalf("Expected log item with ID %d, got %v", id, got_item)
	}

	// data
	data := models.LogData{Data: []byte("Test Data")}
	id, err = logData.Add(data)
	got_data, err := logData.GetByID(id)
	if err != nil || got_data.Id != id {
		t.Fatalf("Expected log data with ID %d, got %v", id, got_data)
	}
}

func ChangeTestHelper(t *testing.T, stp setup) {
	logs, logItems, logData, td := stp(t)
	defer td()

	// log
	log := models.Log{Title: "Test Log", Date: time.Now()}
	id, err := logs.Add(log)
	err = logs.Change(id, models.Log{Id: id, Title: "Changed Log", Date: time.Now()})
	if err != nil {
		t.Fatalf("Failed to change log: %v", err)
	}
	changedLog, err := logs.GetByID(id)
	if err != nil || changedLog.Title != "Changed Log" {
		t.Fatalf("Expected log title 'Changed Log', got %v", changedLog)
	}

	// item
	item := models.LogItem{LogId: id, Type: models.Text, Content: "Test Item", Order: 1}
	id, err = logItems.Add(item)
	err = logItems.Change(id, models.LogItem{Id: id, LogId: id, Type: models.Image, Content: "Changed Item", Order: 2})
	if err != nil {
		t.Fatalf("Failed to change log item: %v", err)
	}
	changedItem, err := logItems.GetByID(id)
	if err != nil || changedItem.Content != "Changed Item" {
		t.Fatalf("Expected log item content 'Changed Item', got %v", changedItem)
	}

	// data
	data := models.LogData{Data: []byte("Test Data")}
	id, err = logData.Add(data)
	err = logData.Change(id, models.LogData{Id: id, Data: []byte("Changed Data")})
	if err != nil {
		t.Fatalf("Failed to change log data: %v", err)
	}
	changedData, err := logData.GetByID(id)
	if err != nil || string(changedData.Data) != "Changed Data" {
		t.Fatalf("Expected log data 'Changed Data', got %v", changedData)
	}
}

func GetByFieldTestHelper(t *testing.T, stp setup) {
	logs, logItems, logData, td := stp(t)
	defer td()

	// log
	log := models.Log{Title: "Test Log", Date: time.Now()}
	id, err := logs.Add(log)
	got_logs, err := logs.GetByField("title", "Test Log")
	if err != nil || id == -1 || len(got_logs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(got_logs))
	}

	// item
	item := models.LogItem{LogId: id, Type: models.Text, Content: "Test Item", Order: 1}
	id, err = logItems.Add(item)
	got_items, err := logItems.GetByField("content", "Test Item")
	if err != nil || id == -1 || len(got_items) != 1 {
		t.Fatalf("Expected 1 log item, got %d", len(got_items))
	}

	// data
	data := models.LogData{Data: []byte("Test Data")}
	id, err = logData.Add(data)
	got_data, err := logData.GetByField("data", []byte("Test Data"))
	if err != nil || id == -1 || len(got_data) != 1 {
		t.Fatalf("Expected 1 log data, got %d", len(got_data))
	}
}

func RemoveTestHelper(t *testing.T, stp setup) {
	logs, logItems, logData, td := stp(t)
	defer td()

	// log
	log := models.Log{Title: "Test Log", Date: time.Now()}
	id, err := logs.Add(log)
	err = logs.Remove(id)
	if err != nil {
		t.Fatalf("Failed to remove log: %v", err)
	}
	_, err = logs.GetByID(id)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	// item
	item := models.LogItem{LogId: id, Type: models.Text, Content: "Test Item", Order: 1}
	id, err = logItems.Add(item)
	err = logItems.Remove(id)
	if err != nil {
		t.Fatalf("Failed to remove log item: %v", err)
	}
	_, err = logItems.GetByID(id)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	// data
	data := models.LogData{Data: []byte("Test Data")}
	id, err = logData.Add(data)
	err = logData.Remove(id)
	if err != nil {
		t.Fatalf("Failed to remove log data: %v", err)
	}
	_, err = logData.GetByID(id)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

// Non trivial helpers

//! Add tests to test non trivial actions (like adding multiple items and getting them by random field and stuff)
