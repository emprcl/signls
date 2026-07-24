package filesystem

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

// Configuration represents a configuration loaded from a json file.
type Configuration struct {
	KeyMap   KeyMap      `json:"keymap"`
	Theme    ThemeConfig `json:"theme"`
	version  string
	filename string
}

// NewConfiguration returns a new default configuration. A non-empty keyboard or
// theme overrides whatever the loaded config file specifies and is persisted
// back to disk.
func NewConfiguration(filename, version, keyboard, theme string) Configuration {
	config := Configuration{
		KeyMap:   NewDefaultQwertyKeyMap(),
		Theme:    ThemeConfig{Name: DefaultThemeName},
		version:  version,
		filename: filename,
	}
	config.Load(filename)

	if keyboard != "" {
		switch keyboard {
		case "qwerty-mac":
			config.KeyMap = NewDefaultQwertyMacKeyMap()
		case "azerty":
			config.KeyMap = NewDefaultAzertyKeyMap()
		case "azerty-mac":
			config.KeyMap = NewDefaultAzertyMacKeyMap()
		}
	}

	// The -theme flag selects a preset by name, leaving any user palette
	// overrides in place (they still merge over the new preset).
	if theme != "" && HasTheme(theme) {
		config.Theme.Name = theme
	}

	// Keep the theme name discoverable in the saved config when the file omits
	// it (e.g. configs written before themes existed).
	if config.Theme.Name == "" {
		config.Theme.Name = DefaultThemeName
	}

	config.Save()

	return config
}

// Palette returns the fully resolved color palette for the configured theme,
// with any user overrides applied on top of the selected preset.
func (c Configuration) Palette() Palette {
	return c.Theme.Resolve()
}

// Version returns the version number.
func (c Configuration) Version() string {
	if c.version == "" {
		return "dev"
	}
	return c.version
}

// Save serializes the Configuration and writes it to a file.
func (c *Configuration) Save() {
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(c.filename, content, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// Load reads a json and unmarshal its content to the Configuration.
func (c *Configuration) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	err = json.Unmarshal(content, c)
	if err != nil {
		log.Fatal(err)
	}
	c.filename = filename
}
