package field

import (
	"cykl/core/music"
	"cykl/filesystem"
	"cykl/midi"
)

func NewFromBank(grid filesystem.Grid, midi midi.Midi) *Grid {
	newGrid := NewGrid(grid.Width, grid.Height, midi)
	newGrid.Load(grid)
	return newGrid
}

func (g *Grid) Save(bank *filesystem.Bank) {
	nodes := []filesystem.Node{}

	for y := range g.nodes {
		for _, n := range g.nodes[y] {
			if n == nil {
				continue
			}

			node := filesystem.Node{}
			node.Type = n.Name()
			node.Direction = int(n.Direction())
			node.Note = filesystem.NewNote(*n.(music.Audible).Note())

			switch node.Type {
			case "bang":

			}

			nodes = append(nodes)
		}
	}

	bank.Save(filesystem.Grid{
		Nodes:  nodes,
		Tempo:  g.Tempo(),
		Height: g.Height,
		Width:  g.Width,
		Key:    uint8(g.Key),
		Scale:  uint16(g.Scale),
	})
}

func (g *Grid) Load(grid filesystem.Grid) {
	g.Reset()
	g.midi.SilenceAll()

	g.mu.Lock()
	defer g.mu.Unlock()

	g.clock.SetTempo(grid.Tempo)
	g.Key = music.Key(grid.Key)
	g.Scale = music.Scale(grid.Scale)
}
