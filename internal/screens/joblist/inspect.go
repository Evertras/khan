package joblist

import tea "github.com/charmbracelet/bubbletea"

func (m Model) updateInspect() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) viewInspect() string {
	return "Inspect"
}
