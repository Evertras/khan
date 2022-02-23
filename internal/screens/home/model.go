package home

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/datatree"
)

type SampleStruct struct {
	Inner struct {
		ID         string
		Another    int
		unexported float64
	}

	Name string
	shh  int

	Job api.Job

	Nums []int
}

type Model struct {
	tree datatree.Model
}

func NewModel() Model {
	sample := SampleStruct{
		Name: "Hello",
		Nums: []int{3, 4, 1},
	}

	sample.Inner.Another = 3

	sample.Inner.ID = "some-id"

	return Model{
		tree: datatree.New(&sample),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	body := strings.Builder{}
	body.WriteString(" Welcome to Khan!  Press the first letter (with shift) of the tabs above to visit each tab.\n\n")
	body.WriteString(" Press 'q' or ctrl+C at any time to quit.")

	body.WriteString("\n")
	body.WriteString(m.tree.View())

	return body.String()
}
