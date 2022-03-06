package list

import (
	"time"

	"github.com/evertras/bubble-table/table"
	"github.com/evertras/khan/internal/styles"
)

func genListTable() table.Model {
	columns := []table.Column{
		table.NewColumn(tableKeyID, "ID", 15),
		table.NewColumn(tableKeyName, "Name", 20),
		table.NewColumn(tableKeyStatus, "Status", 10),
		table.NewColumn(tableKeySubmitDate, "Submit Timestamp", 20),
	}

	return table.New(columns).
		SelectableRows(true).
		Focused(true).
		HeaderStyle(styles.Bold).
		WithPageSize(10).
		SortByDesc(tableKeySubmitDate)
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
			tableKeyID:         job.ID,
			tableKeyName:       job.Name,
			tableKeyStatus:     job.Status,
			tableKeySubmitDate: time.Unix(0, job.SubmitTime).Format("2006-01-02 15:04:05"),
		})

		switch job.Status {
		case "running":
			row.Style = styles.Good

		case "pending":
			row.Style = styles.Warning

		default:
			row.Style = styles.Error
		}

		rows = append(rows, row)
	}

	return rows
}
