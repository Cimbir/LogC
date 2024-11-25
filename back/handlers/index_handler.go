package handlers

import (
	"LogC/back/code"
	"LogC/back/models"
	"html/template"
	"net/http"
	"path/filepath"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Define the data to pass to the template
		data := struct {
			Logs []models.Log
		}{
			Logs: code.ReadLogsFromFile(),
		}

		// Reverse the order of logs
		for i, j := 0, len(data.Logs)-1; i < j; i, j = i+1, j-1 {
			data.Logs[i], data.Logs[j] = data.Logs[j], data.Logs[i]
		}

		// Parse and execute the template
		tmplPath := filepath.Join("front", "index.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// Handle POST request
	} else {
	}
}
