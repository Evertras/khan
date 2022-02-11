package table

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	headers     []Header
	headerStyle lipgloss.Style

	rows []Row
}

func New(headers []Header) Model {
	m := Model{
		headers: make([]Header, len(headers)),
	}

	// Do a full deep copy to avoid unexpected edits
	copy(m.headers, headers)
	for i, header := range m.headers {
		m.headers[i].Style = header.Style.Copy()
	}

	return m
}

func (m Model) HeaderStyle(style lipgloss.Style) Model {
	m.headerStyle = style.Copy()
	return m
}

func (m Model) WithRows(rows []Row) Model {
	m.rows = rows
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	body := strings.Builder{}

	headerStrings := []string{}

	for i, header := range m.headers {
		fmtString := fmt.Sprintf("%%%ds", header.Width)
		headerSection := fmt.Sprintf(fmtString, header.Title)
		borderStyle := lipgloss.NewStyle()

		if i == 0 {
			//borderStyle = borderStyle.BorderStyle(borderHeaderFirst)
			borderStyle = borderStyle.BorderStyle(borderHeaderTriangleFirst)
		} else if i < len(m.headers)-1 {
			borderStyle = borderStyle.BorderStyle(borderHeaderMiddle).BorderTop(true).BorderBottom(true).BorderRight(true)
		} else {
			//borderStyle = borderStyle.BorderStyle(borderHeaderLast).BorderTop(true).BorderBottom(true).BorderRight(true)
			borderStyle = borderStyle.BorderStyle(borderHeaderTriangleLast).BorderTop(true).BorderBottom(true).BorderRight(true)
		}

		headerStrings = append(headerStrings, borderStyle.Render(header.Style.Render(headerSection)))
	}

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...))

	body.WriteString("\n")
	return body.String()
}
