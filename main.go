package main

import (
	"LogC/back/handlers"
	"fmt"
	"net/http"
)

func main() {
	// Registering the handler functions
	http.HandleFunc("/", handlers.IndexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/static"))))
	//http.Handle("/", http.FileServer(http.Dir("front")))
	fmt.Println("Server is running at 8090 port.")
	http.ListenAndServe(":8090", nil)
}
