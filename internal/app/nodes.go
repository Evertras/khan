package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type NodesModel struct {
}

func (m NodesModel) Init() tea.Cmd {
	return nil
}

func (m NodesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "esc":
			menu := MainMenuModel{}
			cmd := menu.Init()
			return menu, cmd
		}
	}

	return m, nil
}

func (m NodesModel) View() string {
	body := strings.Builder{}

	body.WriteString(header("Nodes", "View current nodes in the Nomad cluster"))

	return body.String()
}
