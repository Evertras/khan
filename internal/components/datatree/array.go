package datatree

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDataNodeArray(data reflect.Value, indentLevel int) string {
	result := strings.Builder{}

	elemType := data.Type().Elem()

	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	switch elemType.Kind() {
	case reflect.Struct, reflect.Array, reflect.Slice:
		result.WriteString("\n")
		for i := 0; i < data.Len(); i++ {
			style := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				MarginLeft((indentLevel + 1) * m.indentSize).
				PaddingLeft(1).
				PaddingRight(1)
			entryStr := m.renderDataNode(data.Index(i), 0)
			result.WriteString(style.Render(trimNewline(entryStr)))
			result.WriteString("\n")
		}

	default:
		result.WriteString(fmt.Sprintf("%v", data))
	}

	return result.String()
}
