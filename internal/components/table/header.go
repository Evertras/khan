package table

import "github.com/charmbracelet/lipgloss"

type Header struct {
	Title string
	Key   string
	Width int
	Style lipgloss.Style
}

func NewHeader(key, title string, width int) Header {
	return Header{
		Key:   key,
		Title: title,
		Width: width,
	}
}

func (h Header) WithStyle(style lipgloss.Style) Header {
	h.Style = style.Copy()
	return h
}

var borderHeaderFirst = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "┃",
	Right:       "┃",
	TopRight:    "┳",
	TopLeft:     "┏",
	BottomRight: "┻",
	BottomLeft:  "┗",
}

var borderHeaderTriangleFirst = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "┃",
	Right:       "┃",
	TopRight:    "┳",
	TopLeft:     "◤",
	BottomRight: "┻",
	BottomLeft:  "◣",
}

var borderHeaderMiddle = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "",
	Right:       "┃",
	TopRight:    "┳",
	TopLeft:     "",
	BottomRight: "┻",
	BottomLeft:  "",
}

var borderHeaderLast = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "",
	Right:       "┃",
	TopRight:    "┓",
	TopLeft:     "",
	BottomRight: "┛",
	BottomLeft:  "",
}

var borderHeaderTriangleLast = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "",
	Right:       "┃",
	TopRight:    "◥",
	TopLeft:     "",
	BottomRight: "◢",
	BottomLeft:  "",
}
