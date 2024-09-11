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
	KeyMap   KeyMap `json:"keymap"`
	filename string
}

// NewConfiguration returns a new default configuration.
func NewConfiguration(filename, keyboard string) Configuration {
	config := Configuration{
		KeyMap:   NewDefaultQwertyKeyMap(),
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
	config.Save()

	return config
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
