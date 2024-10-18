// Package filesystem provides interfaces and serializable structures that
// allows saving/loading grid state to/from json files.
package filesystem

import (
	"cykl/core/common"
	"cykl/core/music"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	defaultTempo               = 120.
	defaultRootKey music.Key   = 60 // Middle C
	defaultScale   music.Scale = music.CHROMATIC
	defaultSize                = 20
	maxGrids                   = 32
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

// IsEmpty returns true if the grid is empty (no nodes).
func (g Grid) IsEmpty() bool {
	return len(g.Nodes) == 0
}

// Node represents a grid node that is json serializable.
type Node struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Note      Note   `json:"note"`
	Type      string `json:"type"`
	Direction int    `json:"direction"`
	Muted     bool   `json:"muted"`

	Params map[string]Param `json:"params"`
}

type Note struct {
	Key         Key   `json:"key"`
	Channel     Param `json:"channel"`
	Velocity    Param `json:"velocity"`
	Length      Param `json:"length"`
	Probability int   `json:"probability"`
}

func NewNote(n music.Note) Note {
	return Note{
		Key:         NewKey(*n.Key),
		Channel:     NewParam(*n.Channel),
		Velocity:    NewParam(*n.Velocity),
		Length:      NewParam(*n.Length),
		Probability: int(n.Probability),
	}
}

type Key struct {
	Key    int
	Amount int
	Silent bool
}

func NewKey(key music.KeyValue) Key {
	return Key{
		Key:    int(key.BaseValue()),
		Amount: key.RandomAmount(),
		Silent: key.IsSilent(),
	}
}

type Param struct {
	Value  int
	Min    int
	Max    int
	Amount int
}

func NewParam[T uint8 | int](p common.ControlValue[T]) Param {
	return Param{
		Value:  int(p.Value()),
		Min:    int(p.Min()),
		Max:    int(p.Max()),
		Amount: int(p.RandomAmount()),
	}
}

func NewParamFromFile[T uint8 | int](param Param) *common.ControlValue[T] {
	value := common.NewControlValue[T](
		T(param.Value),
		T(param.Min),
		T(param.Max),
	)
	value.SetRandomAmount(int(param.Amount))
	return value
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

// Filename returns the bank filename.
func (b *Bank) Filename() string {
	return strings.TrimSuffix(b.filename, filepath.Ext(b.filename))
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
