package nodes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/menu"
	"github.com/evertras/khan/internal/styles"
)

type Model struct {
	nodes []*api.NodeListStub

	nodeMenu menu.Model
}

func NewEmptyModel() Model {
	return Model{}
}

func NewModelWithNodes(nodes []*api.NodeListStub) Model {
	menuItems := []menu.Item{menu.ItemBack}

	return Model{
		nodes:    nodes,
		nodeMenu: menu.NewModel(menuItems),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.nodeMenu, cmd = m.nodeMenu.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case []*api.NodeListStub:
		m = NewModelWithNodes(msg)
	}

	switch m.nodeMenu.Selected() {
	case menu.ItemBack.Name():
		cmds = append(cmds, BackCmd)
	}

	return m, tea.Batch(cmds...)
}

const nodeTableFormatter = " %30s | %10s | %10s\n"
var nodeTableHeader = fmt.Sprintf(nodeTableFormatter, "Name", "Status", "Address")

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(m.nodeMenu.View())
	body.WriteString("\n")

	body.WriteString(nodeTableHeader)

	for _, node := range m.nodes {
		line := fmt.Sprintf(nodeTableFormatter, node.Name, node.Status, node.Address)
		switch node.Status {
		case "ready":
			line = styles.Good.Render(line)
		}
		body.WriteString(line)
	}

	return body.String()
}
