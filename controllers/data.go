package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type FeedResponse struct {
	Feeds []struct {
		OnestopID string `json:"onestop_id"`
		Name      string `json:"name"`
		URL       string `json:"url"`
	} `json:"feeds"`
}

func initKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return "", err
	}

	apiKey := os.Getenv("API_KEY_TRANSITLAND")
	if apiKey == "" {
		return "", fmt.Errorf("TRANSITLAND_API_KEY is missing")
	}
	return apiKey, nil
}

func GetData(w http.ResponseWriter, r *http.Request) {
	apiKey, err := initKey()
	if err != nil {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		log.Println("API key error:", err)
		return
	}

	lat, lon, radius := 43.7, -79.4, 50
	url := fmt.Sprintf(
		"https://transit.land/api/v2/rest/routes?apikey=%s&lat=%f&lon=%f&radius=%d",
		apiKey, lat, lon, radius,
	)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Println("Request creation error:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch data from Transitland", http.StatusBadGateway)
		log.Println("API request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error response from Transitland", resp.StatusCode)
		log.Printf("Transitland API returned status: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
		log.Println("Response body read error:", err)
		return
	}

	var feedResponse FeedResponse
	if err := json.Unmarshal(body, &feedResponse); err != nil {
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		log.Println("JSON decode error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feedResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
	}
}
