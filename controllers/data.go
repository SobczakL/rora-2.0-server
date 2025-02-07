package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type FeedResponse struct {
	Feeds []struct {
		OnestopID string `json:"onestop_id"`
		Name      string `json:"name"`
		URL       string `json:"url"`
	} `json:"feeds"`
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	apiKey := os.Getenv("TRANSITLAND_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		log.Println("TRANSITLAND_API_KEY is missing")
		return
	}

	url := fmt.Sprintf("https://transit.land/api/v2/rest/feeds?apikey=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch data from Transitland", http.StatusBadGateway)
		log.Println("Error making API request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read API response", http.StatusBadGateway)
		log.Println("Error reading response body:", err)
		return
	}

	var feedResponse FeedResponse
	if err := json.Unmarshal(body, &feedResponse); err != nil {
		http.Error(w, "Failed to parse API response", http.StatusBadGateway)
		log.Println("Error decoding JSON:", err)
		return
	}

	json.NewEncoder(w).Encode(feedResponse)
}
