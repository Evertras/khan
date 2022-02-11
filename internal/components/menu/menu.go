package menu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	ItemBack = NewItem("Back", "b", "esc")
)

type Item struct {
	name              string
	shortcutKey       string
	extraShortcutKeys []string
}

type Model struct {
	items    []Item
	selected string
}

func NewModel(items []Item) Model {
	return Model{
		items: items,
	}
}

func NewItem(name string, shortcutKey string, extraShortcutKeys ...string) Item {
	if name == "" {
		panic("menu item name cannot be empty")
	}

	if shortcutKey == "" {
		panic("shortcut key cannot be empty")
	}

	return Item{
		name:              name,
		shortcutKey:       shortcutKey,
		extraShortcutKeys: extraShortcutKeys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		keyStr := t.String()

		for _, item := range m.items {
			if item.shortcutKey == keyStr {
				m.selected = item.name
			}

			for _, extra := range item.extraShortcutKeys {
				if extra == keyStr {
					m.selected = item.name
				}
			}
		}
	}
	return m, nil
}

// Selected returns the selected menu item name, or an empty string if none has
// been selected
func (m Model) Selected() string {
	return m.selected
}

func (m Model) View() string {
	body := strings.Builder{}
	for _, item := range m.items {
		body.WriteString(fmt.Sprintf(" %s) %s\n", item.shortcutKey, item.name))
	}

	return body.String()
}
