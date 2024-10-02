// Package filesystem provides interfaces and serializable structures that
// allows saving/loading grid state to/from json files.
package filesystem

import (
	"cykl/core/music"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

const (
	defaultTempo               = 120.
	defaultRootKey music.Key   = 60 // Middle C
	defaultScale   music.Scale = music.CHROMATIC
	defaultSize                = 20
	maxGrids                   = 64
)

// Bank holds a slice of grids in memory
type Bank struct {
	mu sync.Mutex

	Grids    []Grid `json:"grids"`
	Active   int    `json:"active"`
	filename string
}

// Grid holds a grid in memory
type Grid struct {
	Nodes []Node  `json:"nodes"`
	Tempo float64 `json:"tempo"`

	Height int `json:"height"`
	Width  int `json:"width"`

	Key   uint8  `json:"key"`
	Scale uint16 `json:"scale"`
}

// Node represents a grid node that is json serializable.
type Node struct {
	Note      Note   `json:"note"`
	Type      string `json:"type"`
	Direction int    `json:"direction"`
}

type Note struct {
	Behavior string `json:"behavior"`
	Channel ,
	Key:         defaultKey,
	Velocity:    defaultVelocity,
	Length:      defaultLength,
	Probability: maxProbability,
}

type Params struct {
}

// New creates and loads a new bank from a given file.
func New(filename string) *Bank {
	grids := make([]Grid, maxGrids)
	nodes := []Node{}
	for k := range grids {
		grids[k].Nodes = nodes
		grids[k].Height = defaultSize
		grids[k].Width = defaultSize
		grids[k].Tempo = defaultTempo
		grids[k].Key = uint8(defaultRootKey)
		grids[k].Scale = uint16(defaultScale)
	}
	bank := &Bank{
		filename: filename,
		Grids:    grids,
	}
	bank.Read(filename)
	return bank
}

// ActiveGrid returns the active grid from the bank.
func (b *Bank) ActiveGrid() Grid {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Grids[b.Active]
}

// Save saves a grid to the active slot and writes.
func (b *Bank) Save(grid Grid) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Grids[b.Active] = grid
	b.Write()
}

// Write serializes the Bank and writes it to a file.
func (b *Bank) Write() {
	content, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(b.filename, content, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// Read reads a json and unmarshal its content to the Bank..
func (b *Bank) Read(filename string) {
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
