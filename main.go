package main

import (
	"flag"
	"log"

	"cykl/core/field"
	"cykl/filesystem"
	"cykl/midi"
	"cykl/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configFile := flag.String("config", "config.json", "config file to load or create")
	keyboard := flag.String("keyboard", "", "keyboard layout (qwerty, qwerty-mac, azerty, azerty-mac)")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	config := filesystem.NewConfiguration(*configFile, *keyboard)

	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	grid := field.NewGrid(2, 2, midi)

	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	p := tea.NewProgram(ui.New(config, grid))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
