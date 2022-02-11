package nodes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/table"
	"github.com/evertras/khan/internal/styles"
)

type Model struct {
	nodes []*api.NodeListStub

	table table.Model
}

func NewEmptyModel() Model {
	return Model{}
}

const (
	tableKeyName    = "name"
	tableKeyStatus  = "status"
	tableKeyAddress = "address"
)

func NewModelWithNodes(nodes []*api.NodeListStub) Model {
	headers := []table.Header{
		table.NewHeader(tableKeyName, "Name", 30).WithStyle(styles.Bold),
		table.NewHeader(tableKeyStatus, "Status", 10).WithStyle(styles.Bold),
		table.NewHeader(tableKeyAddress, "Address", 20).WithStyle(styles.Bold),
	}

	rows := []table.Row{}

	for _, node := range nodes {
		row := table.NewRow(table.RowData{
			tableKeyName:    node.Name,
			tableKeyStatus:  node.Status,
			tableKeyAddress: node.Address,
		})

		switch node.Status {
		case "ready":
			row.Style = styles.Good

		default:
			row.Style = styles.Error
		}

		rows = append(rows, row)
	}

	return Model{
		nodes: nodes,
		table: table.New(headers).WithRows(rows),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case []*api.NodeListStub:
		m = NewModelWithNodes(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(m.table.View())

	return body.String()
}
