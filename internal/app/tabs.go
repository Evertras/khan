package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/styles"
)

var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(styles.ColorHighlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	tabGapTitle = tabGap.Copy().
			Bold(true).
			Foreground(styles.ColorHighlight)

	tabGapInfo = tabGap.Copy().
			Foreground(styles.ColorSubtle)
)

func (m Model) renderTab(title string, activeWhen activeScreen) string {
	str := styles.Title.Render(title[:1]) + title[1:]
	if activeWhen == m.active {
		return activeTab.Render(str)
	} else {
		return tab.Render(str)
	}
}
