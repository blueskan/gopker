package gopker

import (
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"golang.org/x/net/context"
)

func Containers() ([]types.Container, error) {
	cli, err := docker.NewEnvClient()

	if err != nil {
		return nil, err
	}

	return cli.ContainerList(context.Background(), types.ContainerListOptions{})}

