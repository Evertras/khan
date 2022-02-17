package nodes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/hashicorp/nomad/api"
)

type errMsg error

type Model struct {
	nodes []*api.NodeListStub

	details *api.Node

	table table.Model
}

func NewEmptyModel() Model {
	return Model{
		table: genListTable(),
	}
}

func NewModelWithNodes(nodes []*api.NodeListStub) Model {
	return Model{
		nodes: nodes,
		table: genListTable().WithRows(rowsFromNodes(nodes)),
	}
}

func (m Model) Init() tea.Cmd {
	return refreshNodeListCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Always update received data regardless of view
	switch msg := msg.(type) {
	case []*api.NodeListStub:
		m.nodes = msg
		m.table = m.table.WithRows(rowsFromNodes(msg))

	case *api.Node:
		m.details = msg
	}

	if m.details != nil {
		m, cmd = m.updateDetails(msg)
		cmds = append(cmds, cmd)
	} else {
		m, cmd = m.updateList(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	if m.details != nil {
		body.WriteString(m.viewDetails())
	} else {
		body.WriteString(m.viewList())
	}

	return body.String()
}
