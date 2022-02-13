package joblist

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	lastUpdated time.Time

	confirmStopIDs []string

	errorMessage string
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
		lastUpdated:  time.Now(),
	}

	headers := []table.Column{
		table.NewColumn(tableKeyID, "ID", 15),
		table.NewColumn(tableKeyName, "Name", 20),
		table.NewColumn(tableKeyStatus, "Status", 15),
	}

	rows := m.generateRows()

	m.table = table.New(headers).
		WithRows(rows).
		HeaderStyle(styles.Bold).
		SelectableRows(true).
		Focused(true)

	return m
}

func (m Model) Init() tea.Cmd {
	return refreshJobsCmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.errorMessage = msg.Error()
	}

	if m.errorMessage != "" {
		return m.updateError(msg)
	}

	if len(m.confirmStopIDs) != 0 {
		return m.updateConfirmStop(msg)
	}

	return m.updateMainView(msg)
}

func (m Model) View() string {
	if m.errorMessage != "" {
		return m.viewError()
	}

	if len(m.confirmStopIDs) != 0 {
		return m.viewConfirmStop()
	}

	return m.viewMain()
}
