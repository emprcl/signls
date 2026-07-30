package filesystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Configuration represents a configuration loaded from a json file.
type Configuration struct {
	KeyMap   KeyMap      `json:"keymap"`
	Theme    ThemeConfig `json:"theme"`
	version  string
	filename string
}

// KeyboardNames returns the keyboard layouts the -keyboard flag accepts.
func KeyboardNames() []string {
	return []string{"qwerty", "qwerty-mac", "azerty", "azerty-mac"}
}

// NewConfiguration returns a new default configuration. A non-empty keyboard or
// theme overrides whatever the loaded config file specifies and is persisted
// back to disk.
//
// Every problem it can find (an unknown layout or theme, an unusable color) is
// returned rather than worked around. The ui takes over the screen straight
// after this, so a warning printed here would never be read; and the config is
// saved back on the way out, so quietly ignoring a value the user typed would
// also erase it.
func NewConfiguration(filename, version, keyboard, theme string) (Configuration, error) {
	config := Configuration{
		KeyMap:   NewDefaultQwertyKeyMap(),
		Theme:    ThemeConfig{Name: DefaultThemeName},
		version:  version,
		filename: filename,
	}
	if err := config.Load(filename); err != nil {
		return Configuration{}, fmt.Errorf("%s: %w", filename, err)
	}

	switch keyboard {
	case "": // no flag given: keep the layout the config file holds
	case "qwerty":
		config.KeyMap = NewDefaultQwertyKeyMap()
	case "qwerty-mac":
		config.KeyMap = NewDefaultQwertyMacKeyMap()
	case "azerty":
		config.KeyMap = NewDefaultAzertyKeyMap()
	case "azerty-mac":
		config.KeyMap = NewDefaultAzertyMacKeyMap()
	default:
		return Configuration{}, fmt.Errorf(
			"unknown keyboard layout %q (available layouts: %s)",
			keyboard, strings.Join(KeyboardNames(), ", "),
		)
	}

	// The -theme flag selects a preset by name, leaving any user palette
	// overrides in place (they still merge over the new preset).
	if theme != "" {
		if !HasTheme(theme) {
			return Configuration{}, fmt.Errorf(
				"unknown theme %q (available themes: %s)",
				theme, strings.Join(ThemeNames(), ", "),
			)
		}
		config.Theme.Name = theme
	}

	if err := config.Theme.Validate(); err != nil {
		return Configuration{}, fmt.Errorf("%s: theme: %w", filename, err)
	}

	// Keep the theme name discoverable in the saved config when the file omits
	// it (e.g. configs written before themes existed).
	if config.Theme.Name == "" {
		config.Theme.Name = DefaultThemeName
	}

	config.Save()

	return config, nil
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

// Load reads a json and unmarshal its content to the Configuration. A missing
// file is not an error: it is how a first run starts.
func (c *Configuration) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	if err := json.Unmarshal(content, c); err != nil {
		return err
	}
	c.filename = filename
	return nil
}
