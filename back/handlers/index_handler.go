package handlers

import (
	"LogC/back/code"
	"LogC/back/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Define the data to pass to the template
		data := struct {
			Title  string
			Header string
		}{
			Title:  "LogC",
			Header: "Welcome to LogC",
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
		l := models.Log{
			Date: time.Now(),
			Items: []models.LogItem{
				{
					Type:    models.Text,
					Content: r.PostFormValue("content"),
					Order:   1,
				},
			},
		}

		code.SaveLogsToFile(l)
	} else {
	}
}
