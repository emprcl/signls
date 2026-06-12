package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"signls/core/field"
	"signls/filesystem"
	"signls/midi"
	"signls/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// version is set at build time via -ldflags "-X main.version=...". It defaults
// to "dev" for local builds.
var version = "dev"

func main() {
	configFile := flag.String("config", "", "config file to load or create (default: <user config dir>/signls/config.json)")
	bankFile := flag.String("bank", "", "bank file to store grids (default: <user config dir>/signls/default.json)")
	keyboard := flag.String("keyboard", "", "keyboard layout (qwerty, qwerty-mac, azerty, azerty-mac)")
	showVersion := flag.Bool("version", false, "print current version")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Default config and bank files live in the per-user config directory; an
	// explicit -config/-bank path is honored as given.
	configPath := *configFile
	if configPath == "" {
		configPath = filesystem.DefaultPath("config.json")
	}
	bankPath := *bankFile
	if bankPath == "" {
		bankPath = filesystem.DefaultPath("default.json")
	}

	config := filesystem.NewConfiguration(configPath, version, *keyboard)

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

	bank := filesystem.New(bankPath)
	grid := field.NewFromBank(bank.Active, bank.ActiveGrid(), midi)

	p := tea.NewProgram(ui.New(config, grid, bank))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
