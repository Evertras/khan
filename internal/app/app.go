package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/joblist"
	"github.com/evertras/khan/internal/components/menu"
	"github.com/evertras/khan/internal/components/nodes"
	"github.com/evertras/khan/internal/styles"
)

type activeScreen int

type Model struct {
	mainMenu     menu.Model
	nodesModel   nodes.Model
	joblistModel joblist.Model

	active activeScreen
}

const (
	mainMenuItemNodes = "Nodes"
	mainMenuItemJobs  = "Jobs"
	mainMenuItemQuit  = "Quit"

	activeMainMenu activeScreen = iota
	activeNodes
	activeJobList
)

func NewModel() Model {
	return Model{
		mainMenu: menu.NewModel([]menu.Item{
			menu.NewItem(mainMenuItemNodes, "n"),
			menu.NewItem(mainMenuItemJobs, "j"),
			menu.NewItem(mainMenuItemQuit, "q", "esc"),
		}),

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
		m.mainMenu, cmd = m.mainMenu.Update(msg)
		cmds = append(cmds, cmd)

		switch m.mainMenu.Selected() {
		case mainMenuItemNodes:
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

		case mainMenuItemJobs:
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

		case mainMenuItemQuit:
			cmds = append(cmds, tea.Quit)
		}

	case activeNodes:
		m.nodesModel, cmd = m.nodesModel.Update(msg)
		cmds = append(cmds, cmd)

	case activeJobList:
		m.joblistModel, cmd = m.joblistModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "ctrl+c":
			// Ctrl+C always quits as a safety valve
			cmds = append(cmds, tea.Quit)
		}

	case nodes.BackMsg:
		m.active = activeMainMenu

	case joblist.BackMsg:
		m.active = activeMainMenu
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	switch m.active {
	case activeMainMenu:
		body.WriteString(styles.Header("Khan", "A management tool for Nomad"))
		body.WriteString(m.mainMenu.View())

	case activeNodes:
		body.WriteString(styles.Header("Khan - Nodes", "Current nodes"))
		body.WriteString(m.nodesModel.View())

	case activeJobList:
		body.WriteString(styles.Header("Khan - Jobs", "Current known jobs"))
		body.WriteString(m.joblistModel.View())
	}

	return body.String()
}
