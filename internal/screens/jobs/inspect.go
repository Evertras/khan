package jobs

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateInspect(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.inspect = nil
		}
	}
	return m, nil
}

func (m Model) viewInspect() string {
	return m.inspectDataTree.View()
}
