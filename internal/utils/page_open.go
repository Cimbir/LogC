package utils

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func OpenPage(path string, data any, w http.ResponseWriter) {
	// Parse and execute the template
	templates := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", "header.html"),
		filepath.Join("web", "templates", path),
	))
	err := templates.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
