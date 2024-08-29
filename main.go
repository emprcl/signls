package main

import (
	"flag"
	"log"

	"cykl/core"
	"cykl/midi"
	"cykl/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	grid := core.NewGrid(45, 25, midi)

	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	p := tea.NewProgram(ui.New(grid))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
