package list

import (
	"sort"
	"strings"

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

func genListTable(width int) table.Model {
	// TODO: Configurable somehow, too many things we may want to show
	columns := []table.Column{
		table.NewColumn(tableKeyID, "ID", 10),
		table.NewFlexColumn(tableKeyDatacenter, "Datacenter", 1),
		table.NewFlexColumn(tableKeyName, "Name", 2),
		table.NewColumn(tableKeyStatus, "Status", 8),
		table.NewColumn(tableKeyEligible, "Eligibility", 14),
		table.NewColumn(tableKeyDrain, "Draining", len("Draining")+1),
		table.NewColumn(tableKeyAddress, "Address", 13),
		table.NewFlexColumn(tableKeyDrivers, "Drivers", 2),
		table.NewColumn(tableKeyVersion, "Version", len("Version")+1),
	}

	// Focused can always be true since we only receive updates when this view
	// is active
	return table.New(columns).HeaderStyle(styles.Bold).Focused(true).WithTargetWidth(width)
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
