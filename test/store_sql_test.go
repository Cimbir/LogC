package test

import (
	"LogC/internal/models"
	"LogC/internal/store"
	"database/sql"
	"testing"
)

func SetupSQLTestDB(t *testing.T) (store.DB[models.Log], store.DB[models.LogItem], store.DB[models.LogData], store.DB[models.User], store.DB[models.Comment], teardown) {
	db, err := sql.Open("sqlite3", "test_logs.db")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS test_logs")
	if err != nil {
		t.Fatalf("Failed to drop table test_logs: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS test_log_items")
	if err != nil {
		t.Fatalf("Failed to drop table test_log_items: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS test_log_data")
	if err != nil {
		t.Fatalf("Failed to drop table test_log_data: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS test_users")
	if err != nil {
		t.Fatalf("Failed to drop table test_users: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS test_comments")
	if err != nil {
		t.Fatalf("Failed to drop table test_comments: %v", err)
	}

	store.InitDB("test_logs.db", "test_logs", "test_log_items", "test_log_data", "test_users", "test_comments")

	logs := store.NewSQLDB[models.Log](db, "test_logs")
	logItems := store.NewSQLDB[models.LogItem](db, "test_log_items")
	logData := store.NewSQLDB[models.LogData](db, "test_log_data")
	users := store.NewSQLDB[models.User](db, "test_users")
	comments := store.NewSQLDB[models.Comment](db, "test_comments")

	return logs, logItems, logData, users, comments, func() {
		db.Close()
	}
}

func TestSQLAdd(t *testing.T) {
	AddTestHelper(t, SetupSQLTestDB)
}
func TestSQLGetAll(t *testing.T) {
	GetAllTestHelper(t, SetupSQLTestDB)
}
func TestSQLGetById(t *testing.T) {
	GetByIDTestHelper(t, SetupSQLTestDB)
}
func TestSQLChange(t *testing.T) {
	ChangeTestHelper(t, SetupSQLTestDB)
}
func TestSQLGetByField(t *testing.T) {
	GetByFieldTestHelper(t, SetupSQLTestDB)
}
func TestSQLRemove(t *testing.T) {
	RemoveTestHelper(t, SetupSQLTestDB)
}
