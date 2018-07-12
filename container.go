package gopker

import (
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	dContainer "docker.io/go-docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
)

type containerStatus int

const (
	READY   containerStatus = 0
	RUNNING containerStatus = 1
	STOPPED containerStatus = 2
)

type dockerContext struct {
	DockerApiClient *docker.Client
	DockerContext   *context.Context
}

type portMapping struct {
	hostPort      string
	containerPort string
	protocol      string
}

type container struct {
	ContainerID   string
	Image         string
	Status        containerStatus
	IpAddress     string
	PortMappings  []portMapping
	Volumes       []string
	Environments  []string
	DockerContext *dockerContext
}

func Container(image string) (*container, error) {
	ctx := context.Background()

	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	if r, err := cli.ImagePull(ctx, image, types.ImagePullOptions{}); err != nil {
		return nil, err
	} else {
		io.Copy(os.Stdout, r)
	}

	return &container{
		Image:        image,
		Status:       READY,
		PortMappings: make([]portMapping, 0),
		Volumes:      make([]string, 0),
		DockerContext: &dockerContext{
			DockerApiClient: cli,
			DockerContext:   &ctx,
		}}, nil
}

func (container *container) Port(hostPort string, containerPort string, protocols ...string) *container {
	protocol := "tcp"

	if len(protocols) > 0 {
		protocol = protocols[0]
	} else if len(protocols) > 1 {
		protocol = protocols[0]
		log.Println("Only first protocol is used.")
	}

	container.PortMappings = append(
		container.PortMappings,
		portMapping{hostPort, containerPort, protocol},
	)

	return container
}

func (container *container) Volume(target string) *container {
	container.Volumes = append(container.Volumes, target)

	return container
}

func (container *container) Env(env string) *container {
	container.Environments = append(container.Environments, env)

	return container
}

func (container *container) Start() (*container, error) {
	ctx := container.DockerContext.DockerContext
	cli := container.DockerContext.DockerApiClient

	resp, err := cli.ContainerCreate(*ctx, &dContainer.Config{
		Image: container.Image,
		Env:   container.Environments,
		Tty:   true,
	}, &dContainer.HostConfig{
		Binds:        container.Volumes,
		PortBindings: prepareBindings(container.PortMappings),
	}, nil, "")
	if err != nil {
		return nil, err
	}

	container.ContainerID = resp.ID

	if err := cli.ContainerStart(*ctx, container.ContainerID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	inspectionResult, err := cli.ContainerInspect(*ctx, container.ContainerID)

	if err != nil {
		return nil, err
	}

	container.Status = RUNNING
	container.IpAddress = inspectionResult.NetworkSettings.IPAddress

	return container, nil
}

func (container *container) Stop() error {
	ctx := container.DockerContext.DockerContext
	cli := container.DockerContext.DockerApiClient

	if err := cli.ContainerStop(*ctx, container.ContainerID, nil); err != nil {
		return err
	}

	container.Status = STOPPED

	return nil
}

func prepareBindings(portMappings []portMapping) nat.PortMap {
	portMap := make(nat.PortMap)

	for _, portMapping := range portMappings {
		portBindings := make([]nat.PortBinding, 0)

		portBindings = append(portBindings, nat.PortBinding{"0.0.0.0", portMapping.hostPort})

		containerPort := nat.Port(portMapping.containerPort + "/" + portMapping.protocol)
		portMap[containerPort] = portBindings
	}

	return portMap
}
