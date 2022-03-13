package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/repository"
)

func refreshNodeListCmd() tea.Msg {
	client := repository.GetNomadClient()

	nodes, _, err := client.Nodes().List(&api.QueryOptions{})

	if err != nil {
		return err
	}

	return nodes
}

func detailsSelectedCmd(nodeID string) tea.Cmd {
	return func() tea.Msg {
		return DetailsSelectedID(nodeID)
	}
}
