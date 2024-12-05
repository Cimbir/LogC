package main

import (
	"LogC/internal/handlers"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Change the working directory to the root of the project
	err := os.Chdir(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error changing working directory:", err)
		return
	}

	// Registering the handler functions
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/add", handlers.AddLogHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	fmt.Println("Server is running at 8090 port.")
	http.ListenAndServe(":8090", nil)
}
