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
	tableKeyID         = "id"
	tableKeyName       = "name"
	tableKeyDatacenter = "datacenter"
	tableKeyStatus     = "status"
	tableKeyAddress    = "address"
	tableKeyVersion    = "version"
	tableKeyDrain      = "drain"
	tableKeyEligible   = "eligible"
	tableKeyDrivers    = "drivers"
)

func NewModelWithNodes(nodes []*api.NodeListStub) Model {
	// TODO: Make style global among all headers
	headers := []table.Header{
		table.NewHeader(tableKeyID, "ID", 10).WithStyle(styles.Bold),
		table.NewHeader(tableKeyDatacenter, "Datacenter", 12).WithStyle(styles.Bold),
		table.NewHeader(tableKeyName, "Name", 30).WithStyle(styles.Bold),
		table.NewHeader(tableKeyStatus, "Status", 8).WithStyle(styles.Bold),
		table.NewHeader(tableKeyAddress, "Address", 13).WithStyle(styles.Bold),
		table.NewHeader(tableKeyVersion, "Version", len("Version")+1).WithStyle(styles.Bold),
		table.NewHeader(tableKeyEligible, "Eligible", 14).WithStyle(styles.Bold),
		table.NewHeader(tableKeyDrain, "Draining", len("Draining")+1).WithStyle(styles.Bold),
		table.NewHeader(tableKeyDrivers, "Drivers", 40).WithStyle(styles.Bold),
	}

	rows := []table.Row{}

	for _, node := range nodes {
		data := table.RowData{
			tableKeyID:         node.ID,
			tableKeyDatacenter: node.Datacenter,
			tableKeyName:       node.Name,
			tableKeyStatus:     node.Status,
			tableKeyAddress:    node.Address,
			tableKeyVersion:    node.Version,
			tableKeyEligible:   node.SchedulingEligibility,
			tableKeyDrain:      node.Drain,
		}

		driverStrs := []string{}

		for key := range node.Drivers {
			driverStrs = append(driverStrs, key)
		}

		// TODO: Sort for stability
		data[tableKeyDrivers] = strings.Join(driverStrs, ",")

		row := table.NewRow(data)

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
