package utils

import (
	models "LogC/internal/models/store"
	"LogC/internal/store"
)

type AppData struct {
	Logs       store.DB[models.Log]
	LogItems   store.DB[models.LogItem]
	LogDataCol store.DB[models.LogData]
	Users      store.DB[models.User]
	Comments   store.DB[models.Comment]
}
