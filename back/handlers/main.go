package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/tmp", tmpHandler)
	http.Handle("/", http.FileServer(http.Dir("./front")))
	fmt.Println("Server is running at 8090 port.")
	http.ListenAndServe(":8090", nil)
}

func tmpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Hello, World Get!")
	}
	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "Hello, World Post!")
	}
}
