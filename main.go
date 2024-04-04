package main

import (
	"fmt"
	"sync"
)

func routine() {
	fmt.Println("Started goroutine")

	// Start the server
	Router()
}

func main() {
	fmt.Println("Hello, World!")

	// Start the goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		routine()
		wg.Done()
	}()

	// Wait for the goroutine to finish
	wg.Wait()
}
