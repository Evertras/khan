package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertras/khan/internal/app"
)

func main() {
	p := tea.NewProgram(app.MainMenuModel{})
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
