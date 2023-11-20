package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Global variables to store worker information and maintain state
var (
	workers       []string // List of worker URLs
	workerWeights []int    // Weights of each worker
	currentIndex  int      // Index of the currently selected worker
	balancerMutex sync.Mutex
)

// weightedRoundRobin selects a worker based on weighted round-robin algorithm
func weightedRoundRobin() string {
	balancerMutex.Lock()
	defer balancerMutex.Unlock()

	// Calculate the total weight of all workers
	totalWeight := 0
	for _, weight := range workerWeights {
		totalWeight += weight
	}

	// Choose a worker based on weighted round-robin
	randNum := rand.Intn(totalWeight)
	selectedWorker := ""
	currentWeight := 0
	for i, weight := range workerWeights {
		currentWeight += weight
		if randNum < currentWeight {
			selectedWorker = workers[i]
			currentIndex = i
			break
		}
	}

	return selectedWorker
}

// proxyHandler is the HTTP handler that redirects requests to selected workers
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Select a worker using weighted round-robin
	selectedWorker := weightedRoundRobin()
	fmt.Printf("Redirecting to worker: %s\n", selectedWorker)

	// Forward the request to the selected worker
	resp, err := http.Get(selectedWorker + r.URL.String())
	if err != nil {
		http.Error(w, "Error forwarding request to worker", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers and body to the original response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code of the original response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the original response
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response from worker", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize workers and their weights
	workers = []string{"http://localhost:8080", "http://localhost:8081"} // Add worker URLs here
	workerWeights = []int{1, 2}                                          // Adjust worker weights here

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Set up the HTTP handler for proxying requests
	http.HandleFunc("/", proxyHandler)

	// Start the HTTP server on port 8888
	fmt.Println("Load Balancer is running on :8888")
	http.ListenAndServe(":8888", nil)
}
