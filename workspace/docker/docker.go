package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	imageRef   string = os.Getenv("IMAGE_URL")
	imageUser  string = os.Getenv("IMAGE_USER")
	imageCmd   string = os.Getenv("IMAGE_CMD")
	cli *client.Client
)

func Setup() error {
	var err error
	cli, err = client.NewClientWithOpts()
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return fmt.Errorf("failed to copy stdout: %w", err)
	}

	return nil
}
