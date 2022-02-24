package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/screens/home"
	"github.com/evertras/khan/internal/screens/joblist"
	"github.com/evertras/khan/internal/screens/nodes"
	"github.com/hashicorp/nomad/api"
)

type Model struct {
	screen tea.Model

	activeTab currentActiveTab

	width  int
	height int

	connectionInfo string

	size screens.Size
}

func NewModel() Model {
	return Model{
		screen:         home.NewModel(),
		connectionInfo: api.DefaultConfig().Address,
	}
}

type errMsg struct {
	err error
}

func screenSizeFromWindowSize(msg tea.WindowSizeMsg) screens.Size {
	return screens.Size{
		Width:  msg.Width,
		Height: msg.Height - 3,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	m.screen, cmd = m.screen.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Ctrl+C or q always quit as a safety valve
			cmds = append(cmds, tea.Quit)

		case "H":
			m.screen = home.NewModel()
			m.activeTab = activeHome
			cmds = append(cmds, m.screen.Init())

		case "N":
			m.screen = nodes.NewEmptyModel(m.size)
			m.activeTab = activeNodes
			cmds = append(cmds, m.screen.Init())

		case "J":
			m.screen = joblist.NewEmptyModel(m.size)
			m.activeTab = activeJobList
			cmds = append(cmds, m.screen.Init())
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.size = screenSizeFromWindowSize(msg)

		m.screen, cmd = m.screen.Update(m.size)
		cmds = append(cmds, cmd)
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
		m.renderTab("Home", activeHome),
		m.renderTab("Nodes", activeNodes),
		m.renderTab("Jobs", activeJobList),
		tabGapInfo.Render(m.connectionInfo),
	)

	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))

	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	body.WriteString(row)
	body.WriteString("\n")

	body.WriteString(m.screen.View())

	return body.String()
}
