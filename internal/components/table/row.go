package table

import "github.com/charmbracelet/lipgloss"

type RowData map[string]interface{}

type Row struct {
	Style lipgloss.Style
	Data  RowData
}

func NewRow(data RowData) Row {
	d := Row{
		Data: make(map[string]interface{}),
	}

	for key, val := range data {
		// Doesn't deep copy val, but close enough for now...
		d.Data[key] = val
	}

	return d
}
