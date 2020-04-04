package releases

import (
	"context"

	"github.com/docker/docker/client"
)

func updateDockerImage() (updated bool) {
	updated = true

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		updated = false
		return
	}
	cli.NegotiateAPIVersion(ctx)

	// TODO: build docker image

	return
}
