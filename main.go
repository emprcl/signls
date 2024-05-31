package main

import (
	"cykl/midi"
	"cykl/sequencer"
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

	seq := sequencer.New(midi)

	p := tea.NewProgram(ui.New(seq))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
