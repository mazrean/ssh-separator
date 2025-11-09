package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

var (
	isLocalImage bool
	imageRef     string
	imageUser    string
	imageCmd     string
	cli          *client.Client
)

type Config struct {
	LocalImage  bool
	ImageName   string
	ImageUser   string
	ImageCmd    string
	CPULimit    float64
	MemoryLimit float64
	PidsLimit   int64
}

func SetConfig(config Config) {
	isLocalImage = config.LocalImage
	imageRef = config.ImageName
	imageUser = config.ImageUser
	imageCmd = config.ImageCmd
	cpuLimit = int64(config.CPULimit * 1e9)
	memoryLimit = int64(config.MemoryLimit * 1e6)
	pidsLimit = config.PidsLimit
}

func Setup() error {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	if !isLocalImage {
		reader, err := cli.ImagePull(ctx, imageRef, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return fmt.Errorf("failed to copy stdout: %w", err)
		}
	}

	return nil
}
