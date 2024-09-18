package main

import (
	"Sports-Score-Streaming-Platform/internal/api"
	"Sports-Score-Streaming-Platform/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	apiKey := "24b85a2f148247ac8fd0caaaca161061"

	// Initialize FootballDataClient
	client := api.NewFootballDataClient(apiKey)

	liveMatches, err := client.FetchMatches()

	if err != nil {
		log.Fatalf("Failed to fetch live matches: %v", err)
	}

	err = sendToIngestionService(liveMatches)
	if err != nil {
		log.Fatalf("Failed to send data to ingestion service: %v", err)
	}

	log.Println("Data sent to ingestion service successfully!")
}

// Function to send the fetched data to the ingestion service
func sendToIngestionService(matches []model.Match) error {
	ingestionServiceURL := "http://localhost:8080/ingest"

	// Convert matches to JSON
	jsonData, err := json.Marshal(matches)
	if err != nil {
		return fmt.Errorf("failed to marshal matches: %v", err)
	}

	// Make the POST request to the ingestion service
	req, err := http.NewRequest("POST", ingestionServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set a timeout and use the http.Client to make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ingestion service returned status code: %d", resp.StatusCode)
	}

	return nil
}
