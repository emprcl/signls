package field

import (
	"log"
	"reflect"

	"signls/core/common"
	"signls/core/music"
	"signls/core/node"
	"signls/core/theory"
	"signls/filesystem"
	"signls/midi"
)

// cell identifies a grid position. Used to match live nodes against their
// serialized counterparts when restoring a document.
type cell struct {
	x, y int
}

func NewFromBank(bankIndex int, grid filesystem.Grid, midi midi.Midi) *Grid {
	newGrid := NewGrid(grid.Width, grid.Height, midi, grid.Device)
	newGrid.Load(bankIndex, grid)
	return newGrid
}

// Serialize returns the serializable representation of the grid, taken under the
// read lock. Unlike Save it doesn't touch the bank or write to disk; the ui uses
// it to snapshot editable state for the undo/redo history. Transient playback
// state (moving signals, pulse) is intentionally excluded by snapshotForSave.
func (g *Grid) Serialize() filesystem.Grid {
	var fsGrid filesystem.Grid
	g.Read(func() {
		fsGrid = g.snapshotForSave()
	})
	return fsGrid
}

// RestoreNodes swaps in the document described by fsGrid — nodes, layout,
// default device and the clock/transport flags — without disturbing playback.
// The ui uses it to apply an undo/redo step to a running sequencer.
//
// Unlike Load it keeps the pulse counter, the playing flag and every travelling
// signal, and it reuses the live object of every node whose serialized form the
// restore doesn't actually change. So a node the step leaves alone keeps its
// pulse, its armed state and its behavior position: undoing an edit mid-set
// doesn't realign euclid phases, cut signals in flight or silence the whole rig.
//
// The grid-level musical state (root key, scale, tempo) is deliberately left
// alone. Meta commands write those from the clock goroutine, so they aren't part
// of the user's document; the ui replays user changes to them separately.
func (g *Grid) RestoreNodes(fsGrid filesystem.Grid) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Index the live nodes by cell, in the same serialized form fsGrid uses, so
	// unchanged cells can be recognized and their live objects reused.
	live := g.nodes
	liveNodes := make(map[cell]filesystem.Node, len(live))
	for _, n := range g.snapshotNodes() {
		liveNodes[cell{n.X, n.Y}] = n
	}

	g.device = g.midi.NewDevice(fsGrid.Device, "")
	g.SendClock = fsGrid.SendClock
	g.SendTransport = fsGrid.SendTransport
	g.Width, g.Height = fsGrid.Width, fsGrid.Height

	kept := make(map[cell]bool, len(liveNodes))
	nodes := g.buildNodes(fsGrid, func(n filesystem.Node) common.Node {
		c := cell{n.X, n.Y}
		if l, ok := liveNodes[c]; !ok || !reflect.DeepEqual(l, n) {
			return nil
		}
		kept[c] = true
		return live[n.Y][n.X]
	})

	for y := range live {
		for x := range live[y] {
			switch n := live[y][x].(type) {
			case *node.Signal:
				// Travelling signals aren't part of the document, so carry them
				// over unless a node now occupies their cell.
				if !g.outOfBounds(x, y) && nodes[y][x] == nil {
					nodes[y][x] = n
				}
			case music.Audible:
				// A dropped emitter may hold a sounding note whose note off
				// would never be sent once the object is gone. Stop just those,
				// rather than the blanket all-channels silence Load does.
				if !kept[cell{x, y}] {
					n.Note().Stop()
				}
			}
		}
	}

	g.nodes = nodes
}

func (g *Grid) Save(bank *filesystem.Bank) {
	// Save runs in a background command goroutine, so read the grid and node
	// state under the read lock to build the serializable snapshot, then write
	// to disk (below) without holding the lock.
	var fsGrid filesystem.Grid
	g.Read(func() {
		fsGrid = g.snapshotForSave()
	})
	bank.Save(fsGrid)
}

// snapshotForSave builds the serializable representation of the grid. The
// caller must hold the lock.
func (g *Grid) snapshotForSave() filesystem.Grid {
	return filesystem.Grid{
		Nodes:         g.snapshotNodes(),
		Tempo:         g.Tempo(),
		Height:        g.Height,
		Width:         g.Width,
		Device:        g.device.Name,
		Key:           uint8(g.Key),
		Scale:         uint16(g.Scale),
		SendClock:     g.SendClock,
		SendTransport: g.SendTransport,
	}
}

// snapshotNodes builds the serializable representation of every node on the
// grid. Transient playback state (moving signals) is intentionally excluded. The
// caller must hold the lock.
func (g *Grid) snapshotNodes() []filesystem.Node {
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

	return nodes
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
	g.Width, g.Height = grid.Width, grid.Height
	g.nodes = g.buildNodes(grid, nil)
}

// buildNodes allocates the node matrix described by fsGrid. g.Width and g.Height
// must already match fsGrid. When reuse is non-nil it is consulted for each
// serialized node first: a node it returns is placed as-is instead of being
// rebuilt, which lets a restore keep the live objects of unchanged cells. The
// caller must hold the write lock.
func (g *Grid) buildNodes(fsGrid filesystem.Grid, reuse func(filesystem.Node) common.Node) [][]common.Node {
	nodes := make([][]common.Node, g.Height)
	for i := range nodes {
		nodes[i] = make([]common.Node, g.Width)
	}

	for _, n := range fsGrid.Nodes {
		if g.outOfBounds(n.X, n.Y) {
			continue
		}
		if reuse != nil {
			if reused := reuse(n); reused != nil {
				nodes[n.Y][n.X] = reused
				continue
			}
		}
		if newNode := g.newNode(n); newNode != nil {
			nodes[n.Y][n.X] = newNode
		}
	}

	return nodes
}

// newNode builds a live node from its serialized form, returning nil for an
// unknown node type. The caller must hold the write lock.
func (g *Grid) newNode(n filesystem.Node) common.Node {
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
		return nil
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
		a.Note().Device.Device = device
		a.Note().Device.Enabled = device.Enabled()

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

	return newNode
}
