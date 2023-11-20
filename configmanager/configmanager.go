package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Config struct represents the structure of the configuration file.
type Config struct {
	Workers []string `json:"workers"`
	Weights []int    `json:"weights"`
}

func main() {
	// Specify the path to the configuration file
	configFile := "C:/Users/Asit/Desktop/Coding/Newgolang/load-balancer/config.json"

	// Read the configuration file
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		log.Fatalf("Error reading configuration file: %s", err.Error())
		return
	}

	// Create a Config instance to hold the configuration data
	var config Config

	// Unmarshal the JSON configuration data into the Config struct
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Println("Error decoding configuration:", err)
		log.Fatalf("Error decoding configuration: %s", err.Error())
		return
	}

	// Spawn Load-Balancer and workers

	// Run the Load-Balancer in a goroutine
	go func() {
		cmd := exec.Command("go", "run", "C:/Users/Asit/Desktop/Coding/Newgolang/load-balancer/loadbalancer/loadbalancer.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running Load-Balancer:", err)
		}
	}()

	// Run each worker in a separate goroutine
	for i, workerURL := range config.Workers {
		go func(i int, workerURL string) {
			cmd := exec.Command("go", "run", "C:/Users/Asit/Desktop/Coding/Newgolang/load-balancer/worker/worker.go")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error running Worker %d: %v\n", i, err)
			}
		}(i, workerURL)
	}

	// Keep the Configuration Manager running
	fmt.Println("Configuration Manager is running. Press Enter to exit.")
	fmt.Scanln()
}
