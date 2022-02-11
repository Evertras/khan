package nodes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/menu"
)

type Model struct {
	nodes []*api.NodeListStub

	nodeMenu menu.Model
}

func NewEmptyModel() Model {
	return Model{}
}

func NewModelWithNodes(nodes []*api.NodeListStub) Model {
	items := []menu.Item{menu.ItemBack}

	shortcuts := []string{
		"1",
		"2",
		"3",
		"4",
	}

	for index, node := range nodes {
		items = append(items, menu.NewItem(node.Name, shortcuts[index]))
	}

	return Model{
		nodes:    nodes,
		nodeMenu: menu.NewModel(items),
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

	switch t := msg.(type) {
	case []*api.NodeListStub:
		m = NewModelWithNodes(t)
	}

	switch m.nodeMenu.Selected() {
	case menu.ItemBack.Name():
		cmds = append(cmds, BackCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.nodeMenu.View()
}
