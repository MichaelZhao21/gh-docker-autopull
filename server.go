package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
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
	// Create start of log message
	start := fmt.Sprintf("%s %s | %s |", req.Method, req.URL.Path, time.Now().UTC().Format("2006-01-02T15:04:05-0700"))

	// Make sure it's a post request
	if req.Method != "POST" {
		fmt.Printf("%s Invalid request method: %s\n", start, req.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Make sure the request body is of type application/json
	if req.Header.Get("Content-Type") != "application/json" {
		fmt.Printf("%s Invalid content type: %s\n", start, req.Header.Get("Content-Type"))
		http.Error(w, "Invalid request content type, must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Parse the request body
	var event PushEvent
	err := json.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		fmt.Printf("%s Error parsing request body: %s\n", start, err.Error())
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}

	// Ignore if no commits or created/deleted is true
	if len(event.Commits) == 0 || event.Created || event.Deleted {
		fmt.Printf("%s No commits or create/delete event\n", start)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Ignore if wrong repo
	if event.Repository.FullName != "MichaelZhao21/test" {
		fmt.Printf("%s Not correct repository: %s\n", start, event.Repository.FullName)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Ignore if wrong branch
	branch, _ := strings.CutPrefix(event.Ref, "refs/heads/")
	if branch != "master" {
		fmt.Printf("%s Not correct branch: %s\n", start, branch)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Pull da shi
	fmt.Printf("%s Received push event for %s on branch %s\n", start, event.Repository.FullName, branch)
}

func Router(port string) {
	// Add routes
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/event", webhookEvent)

	// Start server
	fmt.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
