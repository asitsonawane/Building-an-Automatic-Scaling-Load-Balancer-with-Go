package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

// statsMutex is a mutex to safely access and modify the request count.
var statsMutex sync.Mutex

// requestCount keeps track of the total number of requests.
var requestCount int

// Stats represents the structure for storing statistics.
type Stats struct {
	RequestCount int `json:"request_count"`
}

// updateStats increments the request count and updates the stats in a file.
func updateStats() {
	// Lock the mutex to ensure exclusive access to requestCount.
	statsMutex.Lock()
	defer statsMutex.Unlock()

	// Increment the request count.
	requestCount++

	// Create a Stats instance with the updated count.
	stats := Stats{RequestCount: requestCount}

	// Marshal the stats to JSON format.
	jsonStats, _ := json.Marshal(stats)

	// Write the JSON stats to a file named "stats.json" with read-write permissions.
	_ = os.WriteFile("stats.json", jsonStats, 0644)
}

// helloHandler is an HTTP handler for the "/api/v1/hello" endpoint.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Update the statistics when handling a request.
	updateStats()

	// Prepare a JSON response with a greeting message.
	response := map[string]interface{}{
		"message": "Hello, World!",
	}

	// Encode the JSON response and write it to the response writer.
	json.NewEncoder(w).Encode(response)
}

// statsHandler is an HTTP handler for the "/worker/stats" endpoint.
func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Lock the mutex to ensure exclusive access to stats during reading.
	statsMutex.Lock()
	defer statsMutex.Unlock()

	// Read the stats from the "stats.json" file.
	statsData, err := os.ReadFile("stats.json")
	if err != nil {
		// If there is an error reading the file, return an internal server error.
		http.Error(w, "Error reading stats", http.StatusInternalServerError)
		return
	}

	// Create a Stats instance to decode the JSON data.
	var stats Stats
	if err := json.Unmarshal(statsData, &stats); err != nil {
		// If there is an error decoding the JSON data, return an internal server error.
		http.Error(w, "Error decoding stats", http.StatusInternalServerError)
		return
	}

	// Encode the stats as JSON and write it to the response writer.
	json.NewEncoder(w).Encode(stats)
}

func main() {
	// Register the HTTP handlers for the "/api/v1/hello" and "/worker/stats" endpoints.
	http.HandleFunc("/api/v1/hello", helloHandler)
	http.HandleFunc("/worker/stats", statsHandler)

	// Print a message indicating that the server is running on port 8080.
	fmt.Println("Worker is running on :8080")

	// Start the HTTP server on port 8080.
	http.ListenAndServe(":8080", nil)

	// Configure log output to a file.
	logFile, err := os.Create("worker.log")
	if err != nil {
		// If there is an error creating the log file, print an error message.
		log.Printf("An error occurred: %s", err.Error())
	}
	defer logFile.Close()
}
