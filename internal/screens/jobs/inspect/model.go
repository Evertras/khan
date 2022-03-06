package inspect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/khan/internal/components/datatree"
	"github.com/hashicorp/nomad/api"
)

type Model struct {
	data *api.Job
	tree datatree.Model
}

func New(data *api.Job) Model {
	return Model{
		data: data,
		tree: datatree.New(data),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tree, cmd = m.tree.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.tree.View()
}
