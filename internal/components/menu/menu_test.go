package menu

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestViewIncludesItemNamesAndMainShortcuts(t *testing.T) {
	items := []Item{
		// Shortcuts as capitals to avoid false positive collisions
		NewItem("abc", "C"),
		NewItem("def", "P", "ctrl+x"),
		NewItem("another", "Z", "ctrl+x"),
	}

	model := NewModel(items)

	rendered := model.View()

	for _, item := range items {
		assert.Contains(t, rendered, item.name)
		assert.Contains(t, rendered, item.shortcutKey)
	}
}

func TestUpdate(t *testing.T) {
	genKeyMsg := func(key rune) tea.Msg {
		return tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{key},
		}
	}

	const menuItemFirst = "first"
	const menuItemSecond = "second"

	items := []Item{
		NewItem(menuItemFirst, "a"),
		NewItem(menuItemSecond, "d", "ctrl+x"),
	}

	model := NewModel(items)

	tests := []struct {
		name              string
		msg               tea.Msg
		expectedSelection string
	}{
		{
			name:              "A random unknown message is sent",
			msg:               17,
			expectedSelection: "",
		},
		{
			name:              "The first shortcut key is pressed",
			msg:               genKeyMsg('a'),
			expectedSelection: menuItemFirst,
		},
		{
			name:              "The second shortcut key is pressed",
			msg:               genKeyMsg('d'),
			expectedSelection: menuItemSecond,
		},
		{
			name: "The second item's extra shortcut key is pressed",
			msg: tea.KeyMsg{
				Type: tea.KeyCtrlX,
			},
			expectedSelection: menuItemSecond,
		},
		{
			name:              "A key is pressed that is not a known shortcut",
			msg:               genKeyMsg('x'),
			expectedSelection: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updated, cmd := model.Update(test.msg)

			assert.Nil(t, cmd, "Don't expect any commands")
			assert.Equal(t, test.expectedSelection, updated.Selected(), "Unexpected selection")
		})
	}
}
