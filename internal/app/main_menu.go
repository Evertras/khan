package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MainMenuModel struct {
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "c":
			var c CountdownModel = 3
			cmd := c.Init()
			return c, cmd

		case "n":
			n := NodesModel{}
			cmd := n.Init()
			return n, cmd

		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MainMenuModel) View() string {
	body := strings.Builder{}

	body.WriteString(header("Khan", "A management tool for Nomad"))

	body.WriteString(" n) Nodes\n")
	body.WriteString(" q) Quit\n")

	return body.String()
}
