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
	DEAD    containerStatus = 3
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
	containerID  string
	image        string
	status       containerStatus
	ipAddress    string
	portMappings []portMapping
	volumes      []string
	environments []string
	context      *dockerContext
}

type Container interface {
	Start() (types.ContainerJSON, error)
	Stop() error
	Kill() error
	Port(string, string, ...string) Container
	Env(string) Container
	Volume(string) Container
}

func NewContainer(image string) (Container, error) {
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
		image:        image,
		status:       READY,
		portMappings: make([]portMapping, 0),
		volumes:      make([]string, 0),
		context: &dockerContext{
			DockerApiClient: cli,
			DockerContext:   &ctx,
		}}, nil
}

func (container *container) Port(hostPort string, containerPort string, protocols ...string) Container {
	protocol := "tcp"

	if len(protocols) > 0 {
		protocol = protocols[0]
	} else if len(protocols) > 1 {
		protocol = protocols[0]
		log.Println("Only first protocol is used.")
	}

	container.portMappings = append(
		container.portMappings,
		portMapping{hostPort, containerPort, protocol},
	)

	return container
}

func (container *container) Volume(target string) Container {
	container.volumes = append(container.volumes, target)

	return container
}

func (container *container) Env(env string) Container {
	container.environments = append(container.environments, env)

	return container
}

func (container *container) Start() (types.ContainerJSON, error) {
	ctx := container.context.DockerContext
	cli := container.context.DockerApiClient

	resp, err := cli.ContainerCreate(*ctx, &dContainer.Config{
		Image: container.image,
		Env:   container.environments,
		Tty:   true,
	}, &dContainer.HostConfig{
		Binds:        container.volumes,
		PortBindings: prepareBindings(container.portMappings),
	}, nil, "")
	if err != nil {
		return types.ContainerJSON{}, err
	}

	container.containerID = resp.ID
	if err := cli.ContainerStart(*ctx, container.containerID, types.ContainerStartOptions{}); err != nil {
		return types.ContainerJSON{}, err
	}

	container.status = RUNNING

	inspectionResult, err := cli.ContainerInspect(*ctx, container.containerID)
	if err != nil {
		return types.ContainerJSON{}, err
	}

	container.ipAddress = inspectionResult.NetworkSettings.IPAddress

	return inspectionResult, nil
}

func (container *container) Stop() error {
	ctx := container.context.DockerContext
	cli := container.context.DockerApiClient

	if err := cli.ContainerStop(*ctx, container.containerID, nil); err != nil {
		return err
	}

	container.status = STOPPED

	return nil
}

func (container *container) Kill() error {
	ctx := container.context.DockerContext
	cli := container.context.DockerApiClient

	if err := cli.ContainerKill(*ctx, container.containerID, ""); err != nil {
		return err
	}

	container.status = DEAD

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
