// Package filesystem provides interfaces and serializable structures that
// allows saving/loading grid state to/from json files.
package filesystem

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

// Grid holds a grid in memory
type Grid struct {
	Nodes [][]Node `json:"nodes"`
	Tempo float64  `json:"tempo"`

	Height int `json:"height"`
	Width  int `json:"width"`

	Key   int `json:"key"`
	Scale int `json:"scale"`

	filename string
}

// Node represents a grid node that is json serializable.
type Node struct {
	Note Note `json:"note"`
	Type
	Params
}

type Note struct {
}

// NewBank creates and loads a new bank from a given file.
func NewBank(filename string) Bank {
	bank := Bank{
		filename: filename,
		Patterns: make([]Pattern, maxPatterns),
	}
	bank.Load(filename)
	return bank
}

// Save serializes the Bank and writes it to a file.
func (b *Bank) Save() {
	content, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(b.filename, content, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// Load reads a json and unmarshal its content to the Bank..
func (b *Bank) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	err = json.Unmarshal(content, b)
	if err != nil {
		log.Fatal(err)
	}
}
