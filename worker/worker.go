package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var statsMutex sync.Mutex
var requestCount int

type Stats struct {
	RequestCount int `json:"request_count"`
}

func updateStats() {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	requestCount++
	stats := Stats{RequestCount: requestCount}
	jsonStats, _ := json.Marshal(stats)
	_ = os.WriteFile("stats.json", jsonStats, 0644)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	updateStats()

	response := map[string]interface{}{
		"message": "Hello, World!",
	}
	json.NewEncoder(w).Encode(response)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Read and return the stats from the file
	statsMutex.Lock()
	defer statsMutex.Unlock()

	statsData, err := os.ReadFile("stats.json")
	if err != nil {
		http.Error(w, "Error reading stats", http.StatusInternalServerError)
		return
	}

	var stats Stats
	if err := json.Unmarshal(statsData, &stats); err != nil {
		http.Error(w, "Error decoding stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func main() {
	http.HandleFunc("/api/v1/hello", helloHandler)
	http.HandleFunc("/worker/stats", statsHandler)

	fmt.Println("Worker is running on :8080")
	http.ListenAndServe(":8080", nil)
	// Configure log output to a file
	logFile, err := os.Create("worker.log")
	if err != nil {
		log.Printf("An error occurred: %s", err.Error())
	}
	defer logFile.Close()
}
