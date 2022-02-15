package errview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertras/khan/internal/styles"
)

var (
	styleErrorMessage = styles.Error.Copy().Padding(2).Bold(true)
)

type Model struct {
	message string
}

func NewEmptyModel() Model {
	return Model{}
}

func NewModelWithMessage(msg string) Model {
	return Model{msg}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Active() bool {
	return m.message != ""
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.message = ""
		}
	}

	return m, nil
}

func (m Model) View() string {
	return styleErrorMessage.Render(fmt.Sprintf("Something went wrong!  Press enter to continue...\n\n%s", m.message))
}
