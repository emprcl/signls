package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"signls/core/field"
	"signls/filesystem"
	"signls/midi"
	"signls/ui"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed VERSION
var AppVersion string

func main() {
	configFile := flag.String("config", "config.json", "config file to load or create")
	bankFile := flag.String("bank", "default.json", "bank file to store grids")
	keyboard := flag.String("keyboard", "", "keyboard layout (qwerty, qwerty-mac, azerty, azerty-mac)")
	version := flag.Bool("version", false, "print current version")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *version {
		fmt.Print(AppVersion)
		os.Exit(0)
	}

	config := filesystem.NewConfiguration(*configFile, strings.TrimSuffix(AppVersion, "\n"), *keyboard)

	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	bank := filesystem.New(*bankFile)
	grid := field.NewFromBank(bank.ActiveGrid(), midi)

	p := tea.NewProgram(ui.New(config, grid, bank))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
