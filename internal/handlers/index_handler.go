package handlers

import (
	"LogC/internal/models"
	"LogC/internal/services"
	"LogC/internal/utils"
	"net/http"
	"path/filepath"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		logs, err := services.ReadLogsFromDB()
		if err != nil {
			http.Error(w, "Error reading logs", http.StatusInternalServerError)
			return
		}

		// Define the data to pass to the template
		data := struct {
			Logs []models.Log
		}{
			//Logs: services.ReadLogsFromFile(),
			Logs: logs,
		}

		// Reverse the order of logs
		for i, j := 0, len(data.Logs)-1; i < j; i, j = i+1, j-1 {
			data.Logs[i], data.Logs[j] = data.Logs[j], data.Logs[i]
		}

		// Parse and execute the template
		path := filepath.Join("web", "templates", "index.html")
		utils.OpenPage(path, data, w)
	} else if r.Method == http.MethodPost {
		// Handle POST request
	} else {
	}
}
