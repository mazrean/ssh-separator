package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/mazrean/separated-webshell/api"
	"github.com/mazrean/separated-webshell/repository/badger"
	"github.com/mazrean/separated-webshell/service"
	"github.com/mazrean/separated-webshell/workspace/docker"
)

type Config struct {
	API struct {
		Port int    `env:"API_PORT" default:"3000" help:"API server port"`
		Key  string `env:"API_KEY" required:"" help:"API authentication key (required for API authentication)"`
	} `embed:"" prefix:""`

	SSH struct {
		Port int `env:"SSH_PORT" default:"2222" help:"SSH server port"`
	} `embed:"" prefix:""`

	Docker struct {
		LocalImage  bool    `env:"LOCAL_IMAGE" default:"false" help:"Use local Docker image"`
		ImageName   string  `env:"IMAGE_NAME" default:"ghcr.io/mazrean/ssh-separator-ubuntu:latest" help:"Docker image name"`
		ImageUser   string  `env:"IMAGE_USER" default:"ubuntu" help:"Docker image user"`
		ImageCmd    string  `env:"IMAGE_CMD" default:"/bin/bash" help:"Docker image command"`
		CPULimit    float64 `env:"CPU_LIMIT" default:"1.0" help:"CPU limit for containers (in CPU cores)"`
		MemoryLimit float64 `env:"MEMORY_LIMIT" default:"1024" help:"Memory limit for containers (in MB)"`
	} `embed:"" prefix:""`

	Connection struct {
		MaxGlobal  int64 `env:"MAX_GLOBAL_CONNECTIONS" default:"1000" help:"Maximum total connections"`
		MaxPerUser int64 `env:"MAX_CONNECTIONS_PER_USER" default:"5" help:"Maximum connections per user"`
	} `embed:"" prefix:""`

	RateLimit struct {
		Rate      int `env:"RATE_LIMIT_RATE" default:"5" help:"Rate limit requests per second"`
		Burst     int `env:"RATE_LIMIT_BURST" default:"5" help:"Rate limit burst size"`
		ExpiresIn int `env:"RATE_LIMIT_EXPIRES_IN" default:"60" help:"Rate limit expiration in seconds"`
	} `embed:"" prefix:""`

	Badger struct {
		Dir string `env:"BADGER_DIR" default:"/var/lib/ssh-separator/badger" help:"Badger database directory"`
	} `embed:"" prefix:""`

	Prometheus bool   `env:"PROMETHEUS" default:"true" help:"Enable Prometheus metrics"`
	Welcome    string `env:"WELCOME" help:"Welcome message displayed on SSH connection"`
}

func main() {
	var config Config
	kong.Parse(&config,
		kong.Name("ssh-separator"),
		kong.Description("Separated Web Shell - SSH to Docker container gateway"),
		kong.UsageOnError(),
	)

	// Set Docker configuration
	docker.SetConfig(docker.Config{
		LocalImage:  config.Docker.LocalImage,
		ImageName:   config.Docker.ImageName,
		ImageUser:   config.Docker.ImageUser,
		ImageCmd:    config.Docker.ImageCmd,
		CPULimit:    config.Docker.CPULimit,
		MemoryLimit: config.Docker.MemoryLimit,
	})

	// Create API configuration
	apiConfig := api.Config{
		Prometheus:         config.Prometheus,
		RateLimitRate:      config.RateLimit.Rate,
		RateLimitBurst:     config.RateLimit.Burst,
		RateLimitExpiresIn: config.RateLimit.ExpiresIn,
	}

	server, closeFn, err := InjectServer(
		api.Key(config.API.Key),
		badger.Dir(config.Badger.Dir),
		service.MaxGlobalConnections(config.Connection.MaxGlobal),
		docker.MaxConnectionsPerUser(config.Connection.MaxPerUser),
		apiConfig,
		service.WelcomeMessage(config.Welcome),
	)
	if err != nil {
		panic(err)
	}
	defer closeFn()

	err = docker.Setup()
	if err != nil {
		panic(fmt.Errorf("failed to setup docker: %w", err))
	}

	err = server.Setup.Setup()
	if err != nil {
		panic(fmt.Errorf("failed to setup service: %w", err))
	}

	apiServer := server.API
	ssh := server.SSH

	go func() {
		panic(apiServer.Start(config.API.Port))
	}()

	panic(ssh.Start(config.SSH.Port))
}
