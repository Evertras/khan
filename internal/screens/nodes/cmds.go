package nodes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/repository"
)

func refreshNodeListCmd() tea.Msg {
	client := repository.GetNomadClient()

	nodes, _, err := client.Nodes().List(&api.QueryOptions{})

	if err != nil {
		return errMsg(err)
	}

	return nodes
}

func getDetailsCmd(id string) func() tea.Msg {
	return func() tea.Msg {
		client := repository.GetNomadClient()

		node, _, err := client.Nodes().Info(id, &api.QueryOptions{})

		if err != nil {
			return errMsg(err)
		}

		return node
	}
}
