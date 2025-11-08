package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/workspace/docker"
)

func main() {
	strAPIPort := os.Getenv("API_PORT")
	strSSHPort := os.Getenv("SSH_PORT")
	apiKey, ok := os.LookupEnv("API_KEY")

	// API_KEY環境変数の検証
	if !ok || apiKey == "" {
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

	// Read max total connections from environment variable
	maxTotalConnectionsStr, ok := os.LookupEnv("MAX_TOTAL_CONNECTIONS")
	var maxTotalConnections int32
	if !ok || maxTotalConnectionsStr == "" {
		maxTotalConnections = 100 // Default value
	} else {
		maxTotalConnectionsInt, err := strconv.ParseInt(maxTotalConnectionsStr, 10, 32)
		if err != nil {
			panic(fmt.Errorf("invalid max total connections: %w", err))
		}
		maxTotalConnections = int32(maxTotalConnectionsInt)
	}

	// Create global connection limiter
	connectionLimiter := domain.NewConnectionLimiter(maxTotalConnections)

	server, close, err := InjectServer(apiKey, connectionLimiter)
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
