package test

import (
	"LogC/internal/utils"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupSQLTestAppData(t *testing.T) (utils.AppData, teardown) {
	logs, logItems, logData, users, comments, td := SetupSQLTestDB(t)

	return utils.AppData{Logs: logs, LogItems: logItems, LogDataCol: logData, Users: users, Comments: comments}, td
}

func TestGetLog(t *testing.T) {
	GetLogTestHelper(t, SetupSQLTestAppData)
}

func TestSaveLog(t *testing.T) {
	SaveLogTestHelper(t, SetupSQLTestAppData)
}

func TestGetComments(t *testing.T) {
	GetCommentsTestHelper(t, SetupSQLTestAppData)
}

func TestSaveComment(t *testing.T) {
	SaveCommentTestHelper(t, SetupSQLTestAppData)
}
