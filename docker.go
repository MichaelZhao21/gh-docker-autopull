package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func dockerBuild(config *Config, repoDir string) error {
	setDockerEnvs(config)

	// Run docker compose instead
	if config.IsCompose {
		return dockerComposeBuild(config, repoDir)
	}

	// Run docker image prune
	err := execCmd("docker", "image", "prune", "-af")
	if err != nil {
		return errors.New("Unable to run docker image prune: " + err.Error())
	}

	// Run docker build
	err = execCmdDir(repoDir, "docker", "build", "-t", config.Tag, ".")
	if err != nil {
		return errors.New("Unable to run docker build: " + err.Error())
	}

	// Run docker stop
	err = execCmd("docker", "stop", config.Tag+"-run")
	if err != nil {
		fmt.Println("\tNo previous docker run to stop")
	}

	// Run docker run
	err = execCmd("docker", "run", "--rm", "-d", "-p", config.PortMap, "--name", config.Tag+"-run", config.Tag)
	if err != nil {
		return errors.New("Unable to run new docker instance: " + err.Error())
	}

	return nil
}

func dockerComposeBuild(config *Config, repoDir string) error {
	err := execCmdDir(repoDir, "docker", "compose", "-p", config.Tag, "build")
	if err != nil {
		return errors.New("Unable to build docker compose image: " + err.Error())
	}

	err = execCmdDir(repoDir, "docker", "compose", "-p", config.Tag, "down")
	if err != nil {
		return errors.New("Unable stop previous compose instance: " + err.Error())
	}

	err = execCmdDir(repoDir, "docker", "compose", "-p", config.Tag, "up", "-d")
	if err != nil {
		return errors.New("Unable to run docker compose: " + err.Error())
	}

	return nil
}

func execCmdDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Run()
	return err
}

func execCmd(name string, args ...string) error {
	return execCmdDir("", name, args...)
}
