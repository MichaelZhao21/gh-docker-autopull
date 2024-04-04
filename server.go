package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"hello": "world"}`))
}

type PushEvent struct {
	Ref        string `json:"ref"`
	Commits    []map[string]interface{}
	Created    bool
	Deleted    bool
	Repository struct {
		FullName string `json:"full_name"`
	}
}

func webhookEvent(w http.ResponseWriter, req *http.Request) {
	// Make sure it's a post request
	if req.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Make sure the request body is of type application/json
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid request content type, must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Read the request body
	body := make([]byte, req.ContentLength)
	_, err := req.Body.Read(body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse the request body
	var event PushEvent
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}

	// Ignore if no commits or created/deleted is true
	if len(event.Commits) == 0 || event.Created || event.Deleted {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Ignore if wrong branch
	if event.Ref != "refs/heads/master" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract the branch from the ref branch
	branch := event.Ref[len("refs/heads/"):]

	// Pull da shi
	fmt.Printf("Received push event for %s on branch %s\n", event.Repository.FullName, branch)
}

func Router() {
	// Add routes
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/event", webhookEvent)

	// Start server
	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
