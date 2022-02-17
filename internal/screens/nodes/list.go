package nodes

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/evertras/khan/internal/styles"
	"github.com/hashicorp/nomad/api"
)

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

func genListTable() table.Model {
	// TODO: Configurable somehow, too many things we may want to show
	columns := []table.Column{
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

	// Focused can always be true since we only receive updates when this view
	// is active
	return table.New(columns).HeaderStyle(styles.Bold).Focused(true)
}

func rowsFromNodes(nodes []*api.NodeListStub) []table.Row {
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

	return rows
}

func (m Model) updateList(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			cmds = append(cmds, refreshNodeListCmd)

		case "enter":
			if len(m.nodes) == 0 {
				break
			}

			cmds = append(cmds, getDetailsCmd(m.table.HighlightedRow().Data[tableKeyID].(string)))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) viewList() string {
	body := strings.Builder{}

	body.WriteString(m.table.View())

	return body.String()
}
