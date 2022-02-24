package nodes

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateDetails(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.details = nil
		}
	}

	return m, cmd
}

func (m Model) viewDetails() string {
	return m.detailsDataTree.View()
}
