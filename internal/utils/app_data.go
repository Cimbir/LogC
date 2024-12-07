package utils

import (
	"LogC/internal/models"
	"LogC/internal/store"
)

type AppData struct {
	Logs       store.DB[models.Log]
	LogItems   store.DB[models.LogItem]
	LogDataCol store.DB[models.LogData]
}
