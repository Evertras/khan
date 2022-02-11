package joblist

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/table"
	"github.com/evertras/khan/internal/styles"
)

type Model struct {
	jobs []*api.JobListStub

	table table.Model

	ShowServices bool
	ShowBatch    bool
}

func NewEmptyModel() Model {
	return Model{
		ShowServices: true,
		ShowBatch:    true,
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
		ShowServices: true,
		ShowBatch:    true,
	}

	m.table = m.generateTable()

	return m
}

func (m Model) generateTable() table.Model {
	headers := []table.Header{
		table.NewHeader(tableKeyID, "ID", 20).WithStyle(styles.Bold),
		table.NewHeader(tableKeyName, "Name", 30).WithStyle(styles.Bold),
		table.NewHeader(tableKeyStatus, "Status", 15).WithStyle(styles.Bold),
	}

	rows := []table.Row{}

JOBLOOP:
	for _, job := range m.jobs {
		switch job.Type {
		case "batch":
			if !m.ShowBatch {
				continue JOBLOOP
			}

		case "service":
			if !m.ShowServices {
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

	return table.New(headers).WithRows(rows)
}

func (m Model) Init() tea.Cmd {
	return nil
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
			m.ShowBatch = !m.ShowBatch
			m.table = m.generateTable()

		case "s":
			m.ShowServices = !m.ShowServices
			m.table = m.generateTable()
		}
	}

	m.table, cmd = m.table.Update(cmd)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("Filters: ")
	body.WriteString(styles.Checkbox("Show (s)ervices", m.ShowServices))
	body.WriteString("  ")
	body.WriteString(styles.Checkbox("Show (b)atch jobs", m.ShowBatch))
	body.WriteString("\n")
	body.WriteString(m.table.View())

	return body.String()
}
