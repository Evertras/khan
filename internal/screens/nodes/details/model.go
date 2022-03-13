package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-data-tree/datatree"
	"github.com/evertras/khan/internal/components/screenviewport"
	"github.com/evertras/khan/internal/repository"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/styles"
	"github.com/hashicorp/nomad/api"
)

type Model struct {
	nodeID string
	data   *api.Node
	size   screens.Size

	detailsDataTree datatree.Model
	viewport        screenviewport.Model
}

func New(nodeID string, size screens.Size) Model {
	return Model{
		nodeID:   nodeID,
		size:     size,
		viewport: screenviewport.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			client := repository.GetNomadClient()

			node, _, err := client.Nodes().Info(m.nodeID, &api.QueryOptions{})

			if err != nil {
				return err
			}

			return node
		},
		func() tea.Msg {
			return m.size
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmds = append(cmds, exitCmd)
		}

	case *api.Node:
		m.data = msg
		m.detailsDataTree = datatree.New(msg)
		m.detailsDataTree, cmd = m.detailsDataTree.Update(m.size)
		cmds = append(cmds, cmd)

		header := lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Node details: " + msg.Name)

		m.viewport = m.viewport.SetContent(m.detailsDataTree.View()).WithHeader(header)

	case screens.Size:
		m.detailsDataTree, cmd = m.detailsDataTree.Update(m.size)
		cmds = append(cmds, cmd)
		m.viewport = m.viewport.SetContent(m.detailsDataTree.View())
	}

	m.detailsDataTree, cmd = m.detailsDataTree.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}
