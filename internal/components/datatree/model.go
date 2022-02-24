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
	}).Bold(true).MarginRight(1)
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

func trimNewline(s string) string {
	return strings.TrimSuffix(s, "\n")
}

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
				MarginLeft((indentLevel + 1) * len(m.indentStr)).
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

func (m Model) renderDataNodeStruct(data reflect.Value, indentLevel int) string {
	result := strings.Builder{}
	indent := strings.Repeat(m.indentStr, indentLevel)

	fieldNames := []string{}

	result.WriteString("\n")

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

		result.WriteString(m.renderDataNode(field, indentLevel+1))
		result.WriteString("\n")
		continue

		switch field.Type().Kind() {
		case reflect.Struct:
			result.WriteString("\n")
			result.WriteString(m.renderDataNodeStruct(field, indentLevel+1))

		case reflect.Slice, reflect.Array:
			result.WriteString(m.renderDataNodeArray(field, indentLevel+1))

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

	var result string

	switch data.Kind() {
	case reflect.Struct:
		result = m.renderDataNodeStruct(data, indentLevel)

	case reflect.Array, reflect.Slice:
		result = m.renderDataNodeArray(data, indentLevel)

	default:
		result = fmt.Sprintf("%v", data)
	}

	return trimNewline(result)
}

func (m Model) View() string {
	body := strings.Builder{}

	reflected := reflect.ValueOf(m.data)

	rendered := m.renderDataNode(reflected, 0)

	body.WriteString(strings.TrimSpace(rendered))

	return body.String()
}
