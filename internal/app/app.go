package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/joblist"
	"github.com/evertras/khan/internal/components/nodes"
)

type activeScreen int

type Model struct {
	nodesModel   nodes.Model
	joblistModel joblist.Model

	width  int
	height int

	active activeScreen
}

const (
	activeMainMenu activeScreen = iota
	activeNodes
	activeJobList
)

func NewModel() Model {
	return Model{
		active: activeMainMenu,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch m.active {
	case activeMainMenu:

	case activeNodes:
		m.nodesModel, cmd = m.nodesModel.Update(msg)
		cmds = append(cmds, cmd)

	case activeJobList:
		m.joblistModel, cmd = m.joblistModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Ctrl+C always quits as a safety valve
			cmds = append(cmds, tea.Quit)

		case "H":
			m.active = activeMainMenu

		case "N":
			m.active = activeNodes
			cmds = append(cmds, func() tea.Msg {
				c, err := api.NewClient(api.DefaultConfig())

				if err != nil {
					panic(err)
				}

				nodes, _, err := c.Nodes().List(&api.QueryOptions{})

				if err != nil {
					panic(err)
				}

				return nodes
			})

		case "J":
			m.active = activeJobList
			cmds = append(cmds, func() tea.Msg {
				c, err := api.NewClient(api.DefaultConfig())

				if err != nil {
					panic(err)
				}

				nodes, _, err := c.Jobs().List(&api.QueryOptions{})

				if err != nil {
					panic(err)
				}

				return nodes
			})
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) View() string {
	body := strings.Builder{}

	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tabGapTitle.Render("Khan"),
		m.renderTab("Home", activeMainMenu),
		m.renderTab("Nodes", activeNodes),
		m.renderTab("Jobs", activeJobList),
	)

	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))

	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	body.WriteString(row)
	body.WriteString("\n")

	switch m.active {
	case activeMainMenu:
		// TODO: Proper lipgloss style
		body.WriteString(" Welcome to Khan!  Press the first letter (with shift) of the tabs above to visit each tab.\n\n")
		body.WriteString(" Press 'q' or ctrl+C at any time to quit.")

	case activeNodes:
		body.WriteString(m.nodesModel.View())

	case activeJobList:
		body.WriteString(m.joblistModel.View())
	}

	return body.String()
}
