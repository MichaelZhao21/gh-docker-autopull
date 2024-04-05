package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func routine(config *Config) {
	fmt.Println("Started goroutine")

	// Start the server
	Router(config.Port)
}

type Config struct {
	// Port to run the server on
	Port string
}

func loadEnvs() *Config {
	// Load the port from the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port: port,
	}
}

func main() {
	fmt.Println("Hello, World!")

	// Load envs
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Did not load .env file")
	} else {
		fmt.Println("Loaded .env file")
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
