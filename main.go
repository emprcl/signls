package main

import (
	"cykl/core"
	"cykl/midi"
	"cykl/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	grid := core.NewGrid(45, 25, midi)

	p := tea.NewProgram(ui.New(grid))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
