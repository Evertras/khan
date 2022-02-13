package nodes

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/repository"
	"github.com/evertras/khan/internal/styles"
)

type errMsg error

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
	headers := []table.Column{
		table.NewColumn(tableKeyID, "ID", 10),
		table.NewColumn(tableKeyDatacenter, "Datacenter", 12),
		table.NewColumn(tableKeyName, "Name", 30),
		table.NewColumn(tableKeyStatus, "Status", 8),
		table.NewColumn(tableKeyEligible, "Eligibility", 14),
		table.NewColumn(tableKeyDrain, "Draining", len("Draining")+1),
		table.NewColumn(tableKeyAddress, "Address", 13),
		table.NewColumn(tableKeyDrivers, "Drivers", 40),
		table.NewColumn(tableKeyVersion, "Version", len("Version")+1),
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

		sort.Strings(driverStrs)

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
		table: table.New(headers).WithRows(rows).HeaderStyle(styles.Bold),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		client := repository.GetNomadClient()

		nodes, _, err := client.Nodes().List(&api.QueryOptions{})

		if err != nil {
			return errMsg(err)
		}

		return nodes
	}
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
