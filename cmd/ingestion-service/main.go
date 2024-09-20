package main

import (
	"Sports-Score-Streaming-Platform/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ingest", ingestHandler)
	log.Println("Ingestion service started, listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ingestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the matches data
	var matches []model.NFLMatch
	if err := json.Unmarshal(body, &matches); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	for _, match := range matches {
		fmt.Printf("Received match: %v\n", match)
	}

	// Respond to the API server
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data ingested successfully"))
}
