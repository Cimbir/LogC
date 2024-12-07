package test

import (
	"LogC/internal/models"
	"LogC/internal/store"
	"database/sql"
	"testing"
)

func setupSqlTestDB(t *testing.T) (store.DB[models.Log], store.DB[models.LogItem], store.DB[models.LogData], teardown) {
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

	store.InitDB("test_logs.db", "test_logs", "test_log_items", "test_log_data")

	logs := store.NewSQLDB[models.Log](db, "test_logs")
	logItems := store.NewSQLDB[models.LogItem](db, "test_log_items")
	logData := store.NewSQLDB[models.LogData](db, "test_log_data")

	return logs, logItems, logData, func() {
		db.Close()
	}
}

func TestSQLAdd(t *testing.T) {
	AddTestHelper(t, setupSqlTestDB)
}
func TestSQLGetAll(t *testing.T) {
	GetAllTestHelper(t, setupSqlTestDB)
}
func TestSQLGetById(t *testing.T) {
	GetByIDTestHelper(t, setupSqlTestDB)
}
func TestSQLChange(t *testing.T) {
	ChangeTestHelper(t, setupSqlTestDB)
}
func TestSQLGetByField(t *testing.T) {
	GetByFieldTestHelper(t, setupSqlTestDB)
}
func TestSQLRemove(t *testing.T) {
	RemoveTestHelper(t, setupSqlTestDB)
}
