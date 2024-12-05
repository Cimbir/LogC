package handlers

import (
	"LogC/internal/services"
	"net/http"
	"strconv"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the image ID from the URL
	imageIDStr := r.URL.Query().Get("id")
	if imageIDStr == "" {
		http.Error(w, "Missing image ID", http.StatusBadRequest)
		return
	}

	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Retrieve the image data from the database
	imageData, err := services.GetImageByID(imageID)
	if err != nil {
		http.Error(w, "Error retrieving image", http.StatusInternalServerError)
		return
	}

	// Set the content type to image
	w.Header().Set("Content-Type", "image/jpeg") // Adjust the content type based on your image format
	w.Write(imageData.Data)
}
