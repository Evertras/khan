package repository

import (
	"sync"

	"github.com/hashicorp/nomad/api"
)

var (
	nomadClientOnce sync.Once
	nomadClient     *api.Client
)

func GetNomadClient() *api.Client {
	nomadClientOnce.Do(func() {
		var err error
		nomadClient, err = api.NewClient(api.DefaultConfig())

		if err != nil {
			// TODO: Better way?
			panic(err)
		}
	})

	return nomadClient
}
