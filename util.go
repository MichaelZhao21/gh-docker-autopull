package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func logMsg(req *http.Request, msg string, msg2 string) {
	start := fmt.Sprintf("%s %s | %s |", req.Method, req.URL.Path, time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	secondPart := ""
	if msg2 != "" {
		secondPart = ": " + msg2
	}
	fmt.Printf("%s %s%s\n", start, msg, secondPart)
}

type Config struct {
	Repo      string // in the format: <org/name>/<repo_name>
	Branch    string // Branch to pull from
	Port      string // Port for THIS application
	IsCompose bool   // Set true if using docker compose
	FileName  string // Dockerfile name or docker compose file name (defaults to Docker/docker-compose.yml)
	Tag       string // Tag of the docker image or name of the project for compose
	PortMap   string // Port map for the running docker container (not required for compose)
	DockerEnv string // Env used to run the dockerfile in
}

func loadEnvs() *Config {
	// Load repo
	repo := os.Getenv("REPO")
	if repo == "" {
		log.Fatalln("REPO environmental variable is required")
	}

	// Load repo
	branch := os.Getenv("BRANCH")
	if branch == "" {
		log.Fatalln("BRANCH environmental variable is required")
	}

	// Load the port from the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Load is compose
	isComposeRaw := os.Getenv("IS_COMPOSE")
	isCompose := true
	if isComposeRaw != "true" && isComposeRaw != "TRUE" {
		isCompose = false
	}

	// Load filename
	fileName := os.Getenv("FILE_NAME")
	if fileName == "" {
		if isCompose {
			fileName = "docker-compose.yml"
		} else {
			fileName = "Dockerfile"
		}
	}

	// Load tag
	tag := os.Getenv("TAG")
	if tag == "" {
		tag = "autopull-image"
	}

	// Load port map
	portMap := os.Getenv("PORT_MAP")
	if portMap == "" {
		portMap = "8080:8080"
	}

	// Load docker env
	dockerEnv := os.Getenv("DOCKER_ENV")

	return &Config{
		Repo:      repo,
		Branch:    branch,
		Port:      port,
		IsCompose: isCompose,
		FileName:  fileName,
		Tag:       tag,
		PortMap:   portMap,
		DockerEnv: dockerEnv,
	}
}

func setDockerEnvs(config *Config) {
	lines := strings.Split(config.DockerEnv, "\n")
	for _, line := range lines {
		// Make sure line is not empty and valid env line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.Contains(line, "=") {
			fmt.Println("\t[WARNING] Env line is invalid: " + line)
			continue
		}

		splitLine := strings.Split(line, "=")
		key := splitLine[0]
		val, found := strings.CutPrefix(line, key+"=")
		if !found {
			val = ""
		}
		err := os.Setenv(key, val)
		if err != nil {
			fmt.Println("\t[WARNING] Unable to set environmental variable: " + err.Error())
		}
	}
}
