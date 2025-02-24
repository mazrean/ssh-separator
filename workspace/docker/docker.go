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
	isLocalImage = os.Getenv("LOCAL_IMAGE")
	imageRef     = os.Getenv("IMAGE_NAME")
	imageUser    = os.Getenv("IMAGE_USER")
	imageCmd     = os.Getenv("IMAGE_CMD")
	cli          *client.Client
)

func Setup() error {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	ctx := context.Background()

	if len(isLocalImage) == 0 || isLocalImage == "false" {
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
