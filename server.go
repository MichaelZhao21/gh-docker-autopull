package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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
		CloneUrl string `json:"clone_url"`
	}
}

func webhookEvent(config *Config) func(http.ResponseWriter, *http.Request) {
	// Return in function wrapper for function handler
	return func(w http.ResponseWriter, req *http.Request) {
		// Make sure it's a post request
		if req.Method != "POST" {
			logMsg(req, "Invalid request method", req.Method)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Make sure the request body is of type application/json
		if req.Header.Get("Content-Type") != "application/json" {
			logMsg(req, "Invalid content type", req.Header.Get("Content-Type"))
			http.Error(w, "Invalid request content type, must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		// Parse the request body
		var event PushEvent
		err := json.NewDecoder(req.Body).Decode(&event)
		if err != nil {
			logMsg(req, "Error parsing request body", err.Error())
			http.Error(w, "Error parsing request body", http.StatusInternalServerError)
			return
		}

		// Ignore if no commits or created/deleted is true
		if len(event.Commits) == 0 || event.Created || event.Deleted {
			logMsg(req, "No commits or create/delete event", "")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Ignore if wrong repo
		if event.Repository.FullName != config.Repo {
			logMsg(req, "Not correct repository", event.Repository.FullName)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Ignore if wrong branch
		branch, _ := strings.CutPrefix(event.Ref, "refs/heads/")
		if branch != config.Branch {
			logMsg(req, "Not correct branch", branch)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log the receive push event
		logMsg(req, fmt.Sprintf("Received push event for %s on branch %s, running build:", event.Repository.FullName, branch), "")

		// Create temp dir for pulling
		err = os.MkdirAll("temp", os.ModePerm)
		if err != nil {
			fmt.Printf("\tUnable to make temp dir\n")
			return
		}

		// Pull da shi
		pullDir := "temp/autopull-" + strings.Replace(event.Repository.FullName, "/", "-", -1)
		_, err = git.PlainClone(pullDir, false, &git.CloneOptions{
			URL:           event.Repository.CloneUrl,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
		})
		if err != nil && err.Error() != "repository already exists" {
			fmt.Printf("\tUnable to clone repository: %s\n", err.Error())
			return
		}

		// Start docker build in go routine
		go func() {
			noError := true

			// Docker build
			err = dockerBuild(config, pullDir)
			if err != nil {
				noError = false
				fmt.Printf("\tError building docker image: %s\n", err.Error())
			}

			// Remove temp dir
			err = os.RemoveAll(pullDir)
			if err != nil {
				fmt.Printf("\tUnable to remove temp directory: %s\n", err.Error())
				noError = false
			}

			if noError {
				fmt.Println("\tSuccessfully finished redeployment!")
			}
		}()
	}
}

func Router(config *Config) {
	// Add routes
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/event", webhookEvent(config))

	// Start server
	fmt.Println("Starting server on port " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
