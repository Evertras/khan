package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertras/khan/internal/components/menu"
	"github.com/evertras/khan/internal/styles"
)

type Model struct {
	mainMenu menu.Model
}

const (
	menuItemNodes = "Nodes"
	menuItemQuit  = "Quit"
)

func NewModel() Model {
	return Model{
		mainMenu: menu.NewModel([]menu.Item{
			menu.NewItem(menuItemNodes, "n"),
			menu.NewItem(menuItemQuit, "q", "esc"),
		}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	m.mainMenu, cmd = m.mainMenu.Update(msg)

	cmds = append(cmds, cmd)

	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "ctrl+c":
			// Ctrl+C always quits as a safety valve
			cmds = append(cmds, tea.Quit)
		}
	}

	switch m.mainMenu.Selected() {
	case menuItemQuit:
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(styles.Header("Khan", "A management tool for Nomad"))

	body.WriteString(m.mainMenu.View())

	return body.String()
}
