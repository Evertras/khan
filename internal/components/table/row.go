package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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

var borderRowLeft = lipgloss.Border{
	Left:        "┃",
	Right:       "┃",
	Bottom:      "━",
	BottomLeft:  "┗",
	BottomRight: "┻",
}

var borderRowMiddle = lipgloss.Border{
	Right:       "┃",
	Bottom:      "━",
	BottomRight: "┻",
}

var borderRowLast = lipgloss.Border{
	Right:       "┃",
	Bottom:      "━",
	BottomRight: "┛",
}

func (r Row) render(headers []Header, last bool) string {
	columnStrings := []string{}

	for i, header := range headers {
		var str string
		if entry, exists := r.Data[header.Key]; exists {
			str = fmt.Sprintf("%v", entry)
		}

		borderStyle := lipgloss.NewStyle()

		if i == 0 {
			borderStyle = borderStyle.BorderStyle(borderRowLeft).BorderRight(true).BorderLeft(true)
		} else if i < len(headers)-1 {
			borderStyle = borderStyle.BorderStyle(borderRowMiddle).BorderRight(true)
		} else {
			borderStyle = borderStyle.BorderStyle(borderRowLast).BorderRight(true)
		}

		if last {
			borderStyle = borderStyle.BorderBottom(true)
		}

		dataStr := fmt.Sprintf(header.fmtString, limitStr(str, header.Width))

		columnStrings = append(columnStrings, borderStyle.Render(dataStr))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, columnStrings...)
}
