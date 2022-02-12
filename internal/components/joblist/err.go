package joblist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateError(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.errorMessage = ""
		}
	}

	return m, nil
}

func (m Model) viewError() string {
	return styleErrorMessage.Render(fmt.Sprintf("Something went wrong!  Press escape to continue...\n\n%s", m.errorMessage))
}
