package datatree

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleFieldKey = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#224",
		Dark:  "#b8e",
	}).Bold(true)
)

type Model struct {
	data      interface{}
	indentStr string
	showZero  bool
}

func New(data interface{}) Model {
	model := Model{
		data:      data,
		indentStr: "  ",
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) renderDataNodeStruct(data reflect.Value, indentLevel int) string {
	result := strings.Builder{}
	indent := strings.Repeat(m.indentStr, indentLevel)

	fieldNames := []string{}

	for i := 0; i < data.Type().NumField(); i++ {
		field := data.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		fieldNames = append(fieldNames, field.Name)
	}

	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		field := data.FieldByName(fieldName)

		if !m.showZero && field.IsZero() {
			continue
		}

		for field.Kind() == reflect.Ptr && !field.IsNil() {
			field = field.Elem()
		}

		result.WriteString(indent)
		result.WriteString(styleFieldKey.Render(fieldName + ":"))

		switch field.Type().Kind() {
		case reflect.Struct:
			result.WriteString("\n")
			result.WriteString(m.renderDataNode(field, indentLevel+1))

		case reflect.Ptr:
			result.WriteString(" ")

			if field.IsNil() {
				result.WriteString("nil")
			} else {
				result.WriteString(fmt.Sprintf("%v", field))
			}
			result.WriteString("\n")

		case reflect.Slice, reflect.Array:
			elemType := field.Type().Elem()

			for elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}

			switch elemType.Kind() {

			case reflect.Struct:
				result.WriteString("\n")
				for i := 0; i < field.Len(); i++ {
					style := lipgloss.NewStyle().
						Border(lipgloss.RoundedBorder()).
						MarginLeft((indentLevel + 1) * len(m.indentStr)).
						PaddingLeft(1).
						PaddingRight(1)
					entryStr := m.renderDataNode(field.Index(i), 0)
					result.WriteString(style.Render(strings.TrimSuffix(entryStr, "\n")))
					result.WriteString("\n")
				}

			default:
				result.WriteString(fmt.Sprintf(" %v\n", field))
			}

		default:
			result.WriteString(fmt.Sprintf(" %v\n", field))
		}
	}

	return result.String()
}

func (m Model) renderDataNode(data reflect.Value, indentLevel int) string {
	for data.Kind() == reflect.Ptr {
		if data.IsNil() {
			return "<nil>"
		}

		data = data.Elem()
	}

	switch data.Kind() {
	case reflect.Struct:
		return m.renderDataNodeStruct(data, indentLevel)

	default:
		return fmt.Sprintf("%v", data)
	}
}

func (m Model) View() string {
	body := strings.Builder{}

	reflected := reflect.ValueOf(m.data)

	body.WriteString(m.renderDataNode(reflected, 0))

	return body.String()
}
