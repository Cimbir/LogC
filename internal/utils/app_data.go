package utils

import (
	"LogC/internal/models"
	"LogC/internal/store"
)

type AppData struct {
	logs      *store.DB[models.Log]
	log_items *store.DB[models.LogItem]
	log_data  *store.DB[models.LogData]
}
