package field

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/core/node"
	"cykl/filesystem"
	"cykl/midi"
	"log"
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

			if _, ok := n.(*node.Signal); ok {
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
	g.Resize(grid.Width, grid.Height)

	g.nodes = make([][]common.Node, g.Height)
	for i := range g.nodes {
		g.nodes[i] = make([]common.Node, g.Width)
	}

	for _, n := range grid.Nodes {
		var newNode common.Node
		switch n.Type {
		case "bang":
			newNode = node.NewBangEmitter(g.midi, common.Direction(n.Direction), true)
		case "euclid":
			newNode = node.NewEuclidEmitter(g.midi, common.Direction(n.Direction))
			newNode.(*node.EuclidEmitter).Steps = filesystem.NewParamFromFile[int](n.Params["steps"])
			newNode.(*node.EuclidEmitter).Triggers = filesystem.NewParamFromFile[int](n.Params["triggers"])
			newNode.(*node.EuclidEmitter).Offset = filesystem.NewParamFromFile[int](n.Params["offset"])
		case "pass":
			newNode = node.NewPassEmitter(g.midi, common.Direction(n.Direction))
		case "relay":
			newNode = node.NewRelayEmitter(g.midi, common.Direction(n.Direction))
		case "cycle":
			newNode = node.NewCycleEmitter(g.midi, common.Direction(n.Direction))
		case "dice":
			newNode = node.NewDiceEmitter(g.midi, common.Direction(n.Direction))
		case "toll":
			newNode = node.NewTollEmitter(g.midi, common.Direction(n.Direction))
			newNode.(common.Behavioral).Behavior().(*node.TollEmitter).Threshold = filesystem.NewParamFromFile[int](n.Params["threshold"])
		case "zone":
			newNode = node.NewZoneEmitter(g.midi, common.Direction(n.Direction))
		case "hole":
			newNode = node.NewHoleEmitter(common.Direction(n.Direction), n.X, n.Y, g.Width, g.Height)
			newNode.(*node.HoleEmitter).DestinationX = filesystem.NewParamFromFile[int](n.Params["destinationX"])
			newNode.(*node.HoleEmitter).DestinationY = filesystem.NewParamFromFile[int](n.Params["destinationY"])
		default:
			log.Printf("cannot load node of type %s", n.Type)
			continue
		}

		if a, ok := newNode.(music.Audible); ok {
			a.SetMute(n.Muted)
			a.Note().SetKey(music.Key(n.Note.Key.Key), g.Key)
			a.Note().Key.SetRandomAmount(n.Note.Key.Amount)
			a.Note().Key.SetSilent(n.Note.Key.Silent)
			a.Note().Channel.Set(uint8(n.Note.Channel.Value))
			a.Note().Channel.SetRandomAmount(n.Note.Channel.Amount)
			a.Note().Channel.SetMin(uint8(n.Note.Channel.Min))
			a.Note().Channel.SetMax(uint8(n.Note.Channel.Max))
			a.Note().Velocity.Set(uint8(n.Note.Velocity.Value))
			a.Note().Velocity.SetRandomAmount(n.Note.Velocity.Amount)
			a.Note().Velocity.SetMin(uint8(n.Note.Velocity.Min))
			a.Note().Velocity.SetMax(uint8(n.Note.Velocity.Max))
			a.Note().Length.Set(uint8(n.Note.Length.Value))
			a.Note().Length.SetRandomAmount(n.Note.Length.Amount)
			a.Note().Length.SetMin(uint8(n.Note.Length.Min))
			a.Note().Length.SetMax(uint8(n.Note.Length.Max))
			a.Note().Probability = uint8(n.Note.Probability)
		}

		g.nodes[n.Y][n.X] = newNode
	}
}
