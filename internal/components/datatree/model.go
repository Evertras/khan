package datatree

import (
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleFieldKey = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#224",
		Dark:  "#b8e",
	}).Bold(true).MarginRight(1)
)

type Model struct {
	data       interface{}
	indentSize int
	showZero   bool
}

func New(data interface{}) Model {
	model := Model{
		data:       data,
		indentSize: 2,
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	body := strings.Builder{}

	reflected := reflect.ValueOf(m.data)

	rendered := m.renderDataNode(reflected, 0)

	body.WriteString(strings.TrimSpace(rendered))

	return body.String()
}
