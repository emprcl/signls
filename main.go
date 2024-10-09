package main

import (
	"flag"
	"log"

	"cykl/core/field"
	"cykl/filesystem"
	"cykl/midi"
	"cykl/ui"

	"image/color"

	"github.com/BigJk/crt"
	bubbleadapter "github.com/BigJk/crt/bubbletea"
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

	fonts, err := crt.LoadFaces("./fonts/IosevkaTerm-Extended.ttf", "./fonts/IosevkaTerm-Extended.ttf", "./fonts/IosevkaTerm-Extended.ttf", crt.GetFontDPI(), 16.0)
	if err != nil {
		log.Fatal(err)
	}

	win, _, err := bubbleadapter.Window(1000, 600, fonts, ui.New(config, grid), color.Black, tea.WithAltScreen())
	if err != nil {
		log.Fatal(err)
	}

	if err := win.Run("Cykl"); err != nil {
		log.Fatal(err)
	}
}
