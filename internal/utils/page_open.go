package utils

import (
	"html/template"
	"net/http"
)

func OpenPage(path string, data any, w http.ResponseWriter) {
	// Parse and execute the template
	temp, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
