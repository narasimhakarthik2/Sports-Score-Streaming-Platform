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
	// footballApiKey := "24b85a2f148247ac8fd0caaaca161061"
	nflApiKey := "a97e771d6dmsh0b6ce358abe48c2p10c0e0jsn202aabd0185c"

	// footballClient := api.NewFootballDataClient(footballApiKey)
	NFLClient := api.NewNFLDataClient(nflApiKey)

	nflWeekMatches, err := NFLClient.FetchGamesForCurrentWeek()
	if err != nil {
		log.Fatalf("Failed to fetch NFL matches: %v", err)
	}
	log.Println(nflWeekMatches)
	if err := sendToIngestionService(nflWeekMatches); err != nil {
		log.Printf("Failed to send NFL data to ingestion service: %v", err)
	} else {
		log.Println("NFL data sent to ingestion service successfully!")
	}
}

func sendToIngestionService(matches []model.NFLMatch) error {
	ingestionServiceURL := "http://localhost:8080/ingest"

	jsonData, err := json.Marshal(matches)
	if err != nil {
		return fmt.Errorf("failed to marshal matches: %v", err)
	}

	req, err := http.NewRequest("POST", ingestionServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ingestion service returned status code: %d", resp.StatusCode)
	}

	return nil
}
