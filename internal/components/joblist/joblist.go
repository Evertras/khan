package joblist

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/styles"
)

type errMsg error

type Model struct {
	jobs []*api.JobListStub

	table table.Model

	showServices bool
	showBatch    bool
}

func NewEmptyModel() Model {
	return Model{
		showServices: true,
		showBatch:    true,
	}
}

const (
	tableKeyID     = "id"
	tableKeyName   = "name"
	tableKeyStatus = "status"
)

func NewModelWithJobs(jobs []*api.JobListStub) Model {
	m := Model{
		jobs:         jobs,
		showServices: true,
		showBatch:    true,
	}

	headers := []table.Header{
		table.NewHeader(tableKeyID, "ID", 15),
		table.NewHeader(tableKeyName, "Name", 20),
		table.NewHeader(tableKeyStatus, "Status", 15),
	}

	rows := m.generateRows()

	m.table = table.New(headers).
		WithRows(rows).
		HeaderStyle(styles.Bold).
		SelectableRows(true).
		Focused(true)

	return m
}

func (m Model) generateRows() []table.Row {
	rows := []table.Row{}

JOBLOOP:
	for _, job := range m.jobs {
		switch job.Type {
		case "batch":
			if !m.showBatch {
				continue JOBLOOP
			}

		case "service":
			if !m.showServices {
				continue JOBLOOP
			}
		}
		row := table.NewRow(table.RowData{
			tableKeyID:     job.ID,
			tableKeyName:   job.Name,
			tableKeyStatus: job.Status,
		})

		switch job.Status {
		case "running":
			row.Style = styles.Good

		default:
			row.Style = styles.Error
		}

		rows = append(rows, row)
	}

	return rows
}

func (m Model) Init() tea.Cmd {
	return refreshJobsCmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case []*api.JobListStub:
		m = NewModelWithJobs(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			m.showBatch = !m.showBatch
			m.table = m.table.WithRows(m.generateRows())

		case "s":
			m.showServices = !m.showServices
			m.table = m.table.WithRows(m.generateRows()).Focused(true)

		case "g":
			cmds = append(cmds, garbageCollectCmd)

		case "r":
			cmds = append(cmds, refreshJobsCmd)
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var (
	styleHelp = lipgloss.NewStyle().Width(70).Padding(1).Foreground(styles.ColorSpecial)
	styleSubtle = lipgloss.NewStyle().Foreground(styles.ColorSubtle)
)

func (m Model) genHelpBox() string {
	deleteHelp := "(d)elete selected"
	if len(m.table.SelectedRows()) == 0 {
		deleteHelp = styleSubtle.Render(deleteHelp)
	}

	return styleHelp.Render("Space/enter to select\n" + deleteHelp + "\n(g)arbage collect\n(r)efresh jobs (clears selections)")
}

func (m Model) View() string {
	if len(m.jobs) == 0 {
		return ""
	}

	tableView := m.table.View()

	body := strings.Builder{}

	body.WriteString("Filters: ")
	body.WriteString(styles.Checkbox("Show (s)ervices", m.showServices))
	body.WriteString("  ")
	body.WriteString(styles.Checkbox("Show (b)atch jobs", m.showBatch))
	body.WriteString("\n")
	body.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		tableView,
		m.genHelpBox(),
	))

	return body.String()
}
