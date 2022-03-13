package list

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/evertras/khan/internal/screens"
	"github.com/hashicorp/nomad/api"
)

type Model struct {
	nodes []*api.NodeListStub
	table table.Model
	size  screens.Size
}

func New(size screens.Size) Model {
	return Model{
		table: genListTable(size.Width),
		size:  size,
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

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			cmds = append(cmds, refreshNodeListCmd)

		case "enter":
			if len(m.nodes) == 0 {
				break
			}

			nodeID := m.table.HighlightedRow().Data[tableKeyID].(string)
			cmds = append(cmds, detailsSelectedCmd(nodeID))
		}

	case []*api.NodeListStub:
		m.nodes = msg
		m.table = m.table.WithRows(rowsFromNodes(msg))

	case screens.Size:
		m.size = msg
		m.table = genListTable(msg.Width).WithRows(rowsFromNodes(m.nodes))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(m.table.View())

	return body.String()
}
