package joblist

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/errview"
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

	errorMessage errview.Model
}

var columns = []table.Column{
	table.NewColumn(tableKeyID, "ID", 15),
	table.NewColumn(tableKeyName, "Name", 20),
	table.NewColumn(tableKeyStatus, "Status", 15),
}

func NewEmptyModel() Model {
	return Model{
		showServices: true,
		showBatch:    true,
		table:        table.New(columns).SelectableRows(true),
		errorMessage: errview.NewEmptyModel(),
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
		errorMessage: errview.NewEmptyModel(),
	}

	rows := m.generateRows()

	m.table = table.New(columns).
		WithRows(rows).
		HeaderStyle(styles.Bold).
		SelectableRows(true).
		Focused(true)

	return m
}

func (m Model) Init() tea.Cmd {
	return refreshJobsCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.errorMessage = errview.NewModelWithMessage(msg.Error())
	}

	if m.errorMessage.Active() {
		var cmd tea.Cmd
		m.errorMessage, cmd = m.errorMessage.Update(msg)
		return m, cmd
	}

	if len(m.confirmStopIDs) != 0 {
		return m.updateConfirmStop(msg)
	}

	return m.updateMainView(msg)
}

func (m Model) View() string {
	if m.errorMessage.Active() {
		return m.errorMessage.View()
	}

	if len(m.confirmStopIDs) != 0 {
		return m.viewConfirmStop()
	}

	return m.viewMain()
}
