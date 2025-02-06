package controllers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Example response data
	response := Response{Message: "Hello from Go controller!"}

	// Return the response as JSON
	json.NewEncoder(w).Encode(response)
}
