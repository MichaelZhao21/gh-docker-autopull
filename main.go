package main

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

func routine(config *Config) {
	// Start the server
	Router(config)
}

func main() {
	fmt.Println("Starting Github Docker Autopull...")

	// Load envs
	err := godotenv.Load()
	if err != nil {
		fmt.Println("\tDid not load .env file")
	} else {
		fmt.Println("\tLoaded .env file")
	}

	// Load the config
	config := loadEnvs()

	// Start the goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		routine(config)
		wg.Done()
	}()

	// Wait for the goroutine to finish
	wg.Wait()
}
