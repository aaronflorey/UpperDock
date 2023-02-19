package container

import (
	"context"
	"dockerup8/internal/config"
	"github.com/docker/docker/api/types"
	"log"
)
import "github.com/docker/docker/client"

type Container struct {
	config config.Config
	api    *client.Client
	ctx    context.Context
}

func New(cfg config.Config) (*Container, error) {

	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	container := &Container{
		config: cfg,
		api:    c,
		ctx:    context.Background(),
	}

	version, err := container.CheckVersion()

	if err != nil {
		return nil, err
	}

	log.Println("Connected to docker running version: ", version.Version)

	return container, err
}

func (c *Container) CheckVersion() (types.Version, error) {
	return c.api.ServerVersion(c.ctx)
}
