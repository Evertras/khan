package list

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateConfirmStop(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		ids := m.confirmStopIDs
		// Regardless, get rid of these
		m.confirmStopIDs = nil

		switch msg.String() {
		case "y":
			return m, stopSelectedCmd(ids)

		default:
			return m, nil
		}
	}

	return m, nil
}

func (m Model) viewConfirmStop() string {
	body := strings.Builder{}

	body.WriteString("WARNING: Are you sure you want to stop the following jobs? [y/N]\n\n")

	body.WriteString(strings.Join(m.confirmStopIDs, ", "))

	return styleConfirmWarning.Render(body.String())
}
