package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"
)

type NodeRepository interface {
	GetNodes() ([]*api.NodeListStub, error)
}

type NodesModel struct {
	repo NodeRepository

	nodes []*api.NodeListStub
}

func NewNodesModel(repo NodeRepository) NodesModel {
	return NodesModel{
		repo: repo,
	}
}

func (m NodesModel) Init() tea.Cmd {
	return nil
}

func (m NodesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m NodesModel) View() string {
	body := strings.Builder{}

	if m.nodes == nil {
		body.WriteString("Fetching nodes...")
	} else {
		for _, node := range m.nodes {
			body.WriteString("  - " + node.Name)
		}
	}

	return body.String()
}
