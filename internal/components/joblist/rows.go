package joblist

import (
	"github.com/evertras/bubble-table/table"
	"github.com/evertras/khan/internal/styles"
)

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
