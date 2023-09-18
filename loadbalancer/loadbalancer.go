package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	workers       []string // List of worker URLs
	workerWeights []int    // Weights of each worker
	currentIndex  int      // Index of the currently selected worker
	balancerMutex sync.Mutex
)

func weightedRoundRobin() string {
	balancerMutex.Lock()
	defer balancerMutex.Unlock()

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

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	selectedWorker := weightedRoundRobin()
	fmt.Printf("Redirecting to worker: %s\n", selectedWorker)

	// Forward the request to the selected worker
	resp, err := http.Get(selectedWorker + r.URL.String())
	if err != nil {
		http.Error(w, "Error forwarding request to worker", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers and body
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response from worker", http.StatusInternalServerError)
		return
	}
}

func main() {
	workers = []string{"http://localhost:8080", "http://localhost:8081"} // Add worker URLs here
	workerWeights = []int{1, 2}                                          // Adjust worker weights here

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", proxyHandler)

	fmt.Println("Load Balancer is running on :8888")
	http.ListenAndServe(":8888", nil)
}
