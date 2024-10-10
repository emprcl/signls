package field

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/core/node"
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
		for x, n := range g.nodes[y] {
			if n == nil {
				continue
			}

			note := filesystem.Note{}
			muted := false
			if a, ok := n.(music.Audible); ok {
				note = filesystem.NewNote(*a.Note())
				muted = a.Muted()
			}

			fnode := filesystem.Node{
				X:         x,
				Y:         y,
				Type:      n.Name(),
				Direction: int(n.Direction()),
				Note:      note,
				Muted:     muted,
				Params:    map[string]filesystem.Param{},
			}

			switch fnode.Type {
			case "euclid":
				fnode.Params = map[string]filesystem.Param{
					"steps":    filesystem.NewParam(*n.(*node.EuclidEmitter).Steps),
					"triggers": filesystem.NewParam(*n.(*node.EuclidEmitter).Triggers),
					"offset":   filesystem.NewParam(*n.(*node.EuclidEmitter).Offset),
				}
			case "toll":
				fnode.Params = map[string]filesystem.Param{
					"threshold": filesystem.NewParam(*n.(common.Behavioral).Behavior().(*node.TollEmitter).Threshold),
				}
			case "hole":
				fnode.Params = map[string]filesystem.Param{
					"destinationX": filesystem.NewParam(*n.(*node.HoleEmitter).DestinationX),
					"destinationY": filesystem.NewParam(*n.(*node.HoleEmitter).DestinationY),
				}
			}

			nodes = append(nodes, fnode)
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
