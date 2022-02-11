package repository

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type Nomad struct {
	client *api.Client
}

func NewDefault() (*Nomad, error) {
	c, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, fmt.Errorf("failed to create Nomad client with default config: %w", err)
	}

	return &Nomad{
		client: c,
	}, nil
}

func (n *Nomad) GetNodes() ([]*api.NodeListStub, error) {
	list, _, err := n.client.Nodes().List(&api.QueryOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get node list: %w", err)
	}

	return list, nil
}
