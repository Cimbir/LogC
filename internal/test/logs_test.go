package test

import (
	"LogC/internal/utils"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupSQLTestAppData(t *testing.T) (utils.AppData, teardown) {
	logs, logItems, logData, td := SetupSQLTestDB(t)

	return utils.AppData{Logs: logs, LogItems: logItems, LogDataCol: logData}, td
}

func TestGetLog(t *testing.T) {
	GetDataTestHelper(t, SetupSQLTestAppData)
}

func TestSaveLog(t *testing.T) {
	SaveDataTestHelper(t, SetupSQLTestAppData)
}
