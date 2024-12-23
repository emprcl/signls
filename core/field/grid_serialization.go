package field

import (
	"log"

	"signls/core/common"
	"signls/core/music"
	"signls/core/node"
	"signls/core/theory"
	"signls/filesystem"
	"signls/midi"
)

func NewFromBank(bankIndex int, grid filesystem.Grid, midi midi.Midi) *Grid {
	newGrid := NewGrid(grid.Width, grid.Height, midi, grid.Device)
	newGrid.Load(bankIndex, grid)
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
			device := ""
			if a, ok := n.(music.Audible); ok {
				note = filesystem.NewNote(*a.Note())
				muted = a.Muted()
				device = a.Note().Device.Name()
			}

			fnode := filesystem.Node{
				X:         x,
				Y:         y,
				Type:      n.Name(),
				Direction: int(n.Direction()),
				Note:      note,
				Muted:     muted,
				Device:    device,
				Params:    map[string]filesystem.Param{},
			}

			switch fnode.Type {
			case "euclid":
				fnode.Params = map[string]filesystem.Param{
					"steps":    filesystem.NewParam(*n.(*node.EuclidEmitter).Steps),
					"triggers": filesystem.NewParam(*n.(*node.EuclidEmitter).Triggers),
					"offset":   filesystem.NewParam(*n.(*node.EuclidEmitter).Offset),
				}
			case "cycle", "dice":
				fnode.Params = map[string]filesystem.Param{
					"repeat": filesystem.NewParam(*n.(common.Behavioral).Behavior().(common.Repeatable).Repeat()),
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
		Nodes:         nodes,
		Tempo:         g.Tempo(),
		Height:        g.Height,
		Width:         g.Width,
		Device:        g.device.Name,
		Key:           uint8(g.Key),
		Scale:         uint16(g.Scale),
		SendClock:     g.SendClock,
		SendTransport: g.SendTransport,
	})
}

func (g *Grid) Load(index int, grid filesystem.Grid) {
	g.Reset()
	g.midi.SilenceAll()

	g.mu.Lock()
	defer g.mu.Unlock()

	g.BankIndex = index
	g.device = g.midi.NewDevice(grid.Device, "")
	g.clock.SetTempo(grid.Tempo)
	g.Key = theory.Key(grid.Key)
	g.Scale = theory.Scale(grid.Scale)
	g.SendClock = grid.SendClock
	g.SendTransport = grid.SendTransport
	g.Resize(grid.Width, grid.Height)

	g.nodes = make([][]common.Node, g.Height)
	for i := range g.nodes {
		g.nodes[i] = make([]common.Node, g.Width)
	}

	for _, n := range grid.Nodes {
		var newNode common.Node
		switch n.Type {
		case "bang":
			newNode = node.NewBangEmitter(g.midi, &g.device, common.Direction(n.Direction), true)
		case "euclid":
			newNode = node.NewEuclidEmitter(g.midi, &g.device, common.Direction(n.Direction))
			newNode.(*node.EuclidEmitter).Steps.Set(n.Params["steps"].Value)
			newNode.(*node.EuclidEmitter).Steps.SetRandomAmount(n.Params["steps"].Amount)
			newNode.(*node.EuclidEmitter).Triggers.Set(n.Params["triggers"].Value)
			newNode.(*node.EuclidEmitter).Triggers.SetRandomAmount(n.Params["triggers"].Amount)
			newNode.(*node.EuclidEmitter).Offset.Set(n.Params["offset"].Value)
			newNode.(*node.EuclidEmitter).Offset.SetRandomAmount(n.Params["offset"].Amount)
		case "pass":
			newNode = node.NewPassEmitter(g.midi, &g.device, common.Direction(n.Direction))
		case "spread":
			newNode = node.NewSpreadEmitter(g.midi, &g.device, common.Direction(n.Direction))
		case "cycle":
			newNode = node.NewCycleEmitter(g.midi, &g.device, common.Direction(n.Direction))
			newNode.(common.Behavioral).Behavior().(*node.CycleEmitter).Repeat().Set(n.Params["repeat"].Value)
			newNode.(common.Behavioral).Behavior().(*node.CycleEmitter).Repeat().SetRandomAmount(n.Params["repeat"].Amount)
		case "dice":
			newNode = node.NewDiceEmitter(g.midi, &g.device, common.Direction(n.Direction))
			newNode.(common.Behavioral).Behavior().(*node.DiceEmitter).Repeat().Set(n.Params["repeat"].Value)
			newNode.(common.Behavioral).Behavior().(*node.DiceEmitter).Repeat().SetRandomAmount(n.Params["repeat"].Amount)
		case "toll":
			newNode = node.NewTollEmitter(g.midi, &g.device, common.Direction(n.Direction))
			newNode.(common.Behavioral).Behavior().(*node.TollEmitter).Threshold.Set(n.Params["threshold"].Value)
			newNode.(common.Behavioral).Behavior().(*node.TollEmitter).Threshold.SetRandomAmount(n.Params["threshold"].Amount)
		case "zone":
			newNode = node.NewZoneEmitter(g.midi, &g.device, common.Direction(n.Direction))
		case "hole":
			newNode = node.NewHoleEmitter(common.Direction(n.Direction), n.X, n.Y, g.Width, g.Height)
			newNode.(*node.HoleEmitter).DestinationX.Set(n.Params["destinationX"].Value)
			newNode.(*node.HoleEmitter).DestinationX.SetRandomAmount(n.Params["destinationX"].Amount)
			newNode.(*node.HoleEmitter).DestinationY.Set(n.Params["destinationY"].Value)
			newNode.(*node.HoleEmitter).DestinationY.SetRandomAmount(n.Params["destinationY"].Amount)
		default:
			log.Printf("cannot load node of type %s", n.Type)
			continue
		}

		if a, ok := newNode.(music.Audible); ok {
			a.SetMute(n.Muted)
			a.Note().SetKey(theory.Key(n.Note.Key.Key), g.Key)
			a.Note().Key.SetRandomAmount(n.Note.Key.Amount)
			a.Note().Key.SetSilent(n.Note.Key.Silent)
			a.Note().Channel.Set(uint8(n.Note.Channel.Value))
			a.Note().Channel.SetRandomAmount(n.Note.Channel.Amount)
			a.Note().Velocity.Set(uint8(n.Note.Velocity.Value))
			a.Note().Velocity.SetRandomAmount(n.Note.Velocity.Amount)
			a.Note().Length.Set(uint8(n.Note.Length.Value))
			a.Note().Length.SetRandomAmount(n.Note.Length.Amount)
			a.Note().Probability = uint8(n.Note.Probability)
			device := g.midi.NewDevice(n.Device, g.device.Name)
			a.Note().Device.Device = g.midi.NewDevice(n.Device, g.device.Name)
			a.Note().Device.Enabled = !device.Empty()

			for i, c := range n.Note.Controls {
				a.Note().Controls[i].Type = music.ControlType(c.Type)
				a.Note().Controls[i].Controller = uint8(c.Controller)
				a.Note().Controls[i].Value.Set(uint8(c.Value.Value))
				a.Note().Controls[i].Value.SetRandomAmount(c.Value.Amount)
			}

			for _, c := range a.Note().MetaCommands {
				cmd := n.Note.MetaCommands[c.Name()]
				c.SetActive(cmd.Active)
				c.Value().Set(cmd.Value.Value)
				c.Value().SetRandomAmount(cmd.Value.Amount)
			}
		}

		g.nodes[n.Y][n.X] = newNode
	}
}
