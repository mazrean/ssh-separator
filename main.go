package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mazrean/separated-webshell/workspace/docker"
)

func main() {
	strAPIPort := os.Getenv("API_PORT")
	strSSHPort := os.Getenv("SSH_PORT")
	apiKey := os.Getenv("API_KEY")

	// API_KEY環境変数の検証
	if apiKey == "" {
		panic(fmt.Errorf("API_KEY environment variable is not set or empty. This is required for API authentication"))
	}

	apiPort, err := strconv.Atoi(strAPIPort)
	if err != nil {
		panic(err)
	}

	sshPort, err := strconv.Atoi(strSSHPort)
	if err != nil {
		panic(err)
	}

	server, close, err := InjectServer(apiKey)
	if err != nil {
		panic(err)
	}
	defer close()

	err = docker.Setup()
	if err != nil {
		panic(fmt.Errorf("failed to setup docker: %w", err))
	}

	err = server.Setup.Setup()
	if err != nil {
		panic(fmt.Errorf("failed to setup service: %w", err))
	}

	api := server.API
	ssh := server.SSH

	go func() {
		panic(api.Start(apiPort))
	}()

	panic(ssh.Start(sshPort))
}
