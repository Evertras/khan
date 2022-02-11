package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// A model can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.
type Model int

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m Model) Init() tea.Cmd {
	return tick
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		m -= 1
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick
	}
	return m, nil
}

// Views return a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	return fmt.Sprintf("Hi. This program will exit in %d seconds. To quit sooner press any key.\n", m)
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
