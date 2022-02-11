package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/menu"
	"github.com/evertras/khan/internal/components/nodes"
	"github.com/evertras/khan/internal/styles"
)

type activeScreen int

type Model struct {
	mainMenu  menu.Model
	nodeModel nodes.Model

	active activeScreen
}

const (
	mainMenuItemNodes = "Nodes"
	mainMenuItemQuit  = "Quit"

	activeMainMenu activeScreen = iota
	activeNodes
)

func NewModel() Model {
	return Model{
		mainMenu: menu.NewModel([]menu.Item{
			menu.NewItem(mainMenuItemNodes, "n"),
			menu.NewItem(mainMenuItemQuit, "q", "esc"),
		}),

		active: activeMainMenu,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
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
		case mainMenuItemQuit:
			cmds = append(cmds, tea.Quit)
		}
	case activeNodes:
		m.nodeModel, cmd = m.nodeModel.Update(msg)
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
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(styles.Header("Khan", "A management tool for Nomad"))

	switch m.active {
	case activeMainMenu:
		body.WriteString(m.mainMenu.View())
	case activeNodes:
		body.WriteString(m.nodeModel.View())
	}

	return body.String()
}
