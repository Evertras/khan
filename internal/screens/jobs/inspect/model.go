package inspect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-data-tree/datatree"
	"github.com/evertras/khan/internal/components/screenviewport"
	"github.com/evertras/khan/internal/screens"
	"github.com/hashicorp/nomad/api"
)

type Model struct {
	data *api.Job
	tree datatree.Model

	viewport screenviewport.Model
	size     screens.Size
}

func New(data *api.Job, size screens.Size) Model {
	return Model{
		data:     data,
		tree:     datatree.New(data),
		viewport: screenviewport.New(),
		size:     size,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return m.size
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case screens.Size:
		m.size = msg
		m.tree, cmd = m.tree.Update(msg)
		cmds = append(cmds, cmd)
		m.viewport = m.viewport.SetContent(m.tree.View())
	}

	m.tree, cmd = m.tree.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}
