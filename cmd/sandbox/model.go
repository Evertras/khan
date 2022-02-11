package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/components/table"
)

type Model struct {
	tableModel table.Model
}

func NewModel() Model {
	headers := []table.Header{
		table.NewHeader("id", "ID", 5).WithStyle(lipgloss.NewStyle().Bold(true)),
		table.NewHeader("name", "Name", 10),
		table.NewHeader("description", "Description", 20),
		table.NewHeader("count", "#", 5),
	}

	rows := []table.Row{
		table.NewRow(table.RowData{
			"id":          "abc",
			"name":        "Hello",
			"description": "The first table entry, ever",
			"count":       4,
		}),
	}

	return Model{
		tableModel: table.New(headers).WithRows(rows),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.tableModel.View()
}
