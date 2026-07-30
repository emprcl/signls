package field

import (
	"sync/atomic"
	"testing"

	"signls/core/common"
	"signls/core/music"
	"signls/core/node"
	"signls/core/theory"
	"signls/midi"
)

// recordingMidi counts the blunt silencing calls, so a test can tell a surgical
// note off from an all-channels blast. The counter is atomic: with the grid
// playing, the clock goroutine sends notes concurrently.
type recordingMidi struct {
	*midi.Mock
	silenceAll atomic.Int64
}

func (m *recordingMidi) SilenceAll() { m.silenceAll.Add(1) }

func newRestoreGrid(t *testing.T) (*Grid, *recordingMidi) {
	t.Helper()
	m := &recordingMidi{Mock: &midi.Mock{}}
	return NewGrid(8, 8, m, ""), m
}

func addBang(grid *Grid, x, y int) {
	device := grid.MidiDevice()
	grid.AddNode(node.NewBangEmitter(grid.Midi(), &device, common.RIGHT, true), x, y)
}

// TestRestoreNodesReusesUnchangedNodes is the point of RestoreNodes: a node the
// step doesn't change keeps its live object, and with it the pulse, armed state
// and behavior position that a rebuild would throw away.
func TestRestoreNodesReusesUnchangedNodes(t *testing.T) {
	grid, m := newRestoreGrid(t)
	addBang(grid, 1, 1)
	addBang(grid, 3, 3)

	doc := grid.Serialize()
	unchanged := grid.Node(1, 1)
	edited := grid.Node(3, 3)
	edited.(music.Audible).Note().Velocity.Set(42)

	grid.RestoreNodes(doc)

	if grid.Node(1, 1) != unchanged {
		t.Error("unchanged node was rebuilt; its playback state would be lost")
	}
	if grid.Node(3, 3) == edited {
		t.Error("edited node was not rebuilt")
	}
	if got := grid.Node(3, 3).(music.Audible).Note().Velocity.Value(); got != 100 {
		t.Errorf("restored velocity: got %d, want 100", got)
	}
	if got := m.silenceAll.Load(); got != 0 {
		t.Errorf("restore silenced every channel (%d calls); it must only stop dropped notes", got)
	}
}

// TestRestoreNodesKeepsTravellingSignals checks signals in flight survive. They
// aren't part of the document, so a restore has nothing to say about them.
func TestRestoreNodesKeepsTravellingSignals(t *testing.T) {
	grid, _ := newRestoreGrid(t)
	addBang(grid, 1, 1)
	doc := grid.Serialize()

	signal := node.NewSignal(common.RIGHT, 0)
	grid.AddNode(signal, 5, 5)

	grid.RestoreNodes(doc)

	if grid.Node(5, 5) != common.Node(signal) {
		t.Error("travelling signal was wiped by the restore")
	}
}

// TestRestoreNodesKeepsPlaybackRunning checks the transport and the pulse counter
// come out the other side. Resetting the pulse would realign every euclid, toll
// and cycle pattern on the grid.
func TestRestoreNodesKeepsPlaybackRunning(t *testing.T) {
	grid, m := newRestoreGrid(t)
	addBang(grid, 1, 1)
	doc := grid.Serialize()

	grid.TogglePlay()
	defer grid.TogglePlay()
	for i := 0; i < common.PulsesPerStep*2; i++ {
		grid.Update()
	}

	grid.RestoreNodes(doc)

	var playing bool
	var pulse uint64
	grid.Read(func() {
		playing = grid.Playing
		pulse = grid.Pulse()
	})
	if !playing {
		t.Error("restore stopped playback")
	}
	if pulse == 0 {
		t.Error("restore reset the pulse counter, realigning every pattern")
	}
	if got := m.silenceAll.Load(); got != 0 {
		t.Errorf("restore silenced every channel (%d calls) mid-playback", got)
	}
}

// TestRestoreNodesLeavesMusicalStateAlone covers the meta command case: root key,
// scale and tempo are written by the clock goroutine, so a restore must not
// revert them. The ui replays the user's own changes to them separately.
func TestRestoreNodesLeavesMusicalStateAlone(t *testing.T) {
	grid, _ := newRestoreGrid(t)
	doc := grid.Serialize()

	scale := theory.AllScales()[1]
	grid.SetKey(theory.Key(72))
	grid.SetScale(scale)
	grid.SetTempo(180)

	grid.RestoreNodes(doc)

	var key theory.Key
	var got theory.Scale
	grid.Read(func() {
		key, got = grid.Key, grid.Scale
	})
	if key != theory.Key(72) {
		t.Errorf("root key: got %d, want 72 — the restore reverted a sequencer change", key)
	}
	if got != scale {
		t.Errorf("scale: got %v, want %v — the restore reverted a sequencer change", got, scale)
	}
}

// TestRestoreNodesRestoresLayout checks undoing a shrink brings back both the
// dimensions and the nodes the shrink dropped.
func TestRestoreNodesRestoresLayout(t *testing.T) {
	grid, _ := newRestoreGrid(t)
	addBang(grid, 6, 6)
	doc := grid.Serialize()

	grid.Resize(4, 4)
	if grid.Width != 4 || grid.Height != 4 {
		t.Fatalf("resize: got %dx%d, want 4x4", grid.Width, grid.Height)
	}

	grid.RestoreNodes(doc)

	if grid.Width != 8 || grid.Height != 8 {
		t.Fatalf("restored size: got %dx%d, want 8x8", grid.Width, grid.Height)
	}
	if grid.Node(6, 6) == nil {
		t.Error("node dropped by the shrink was not restored")
	}
}

// TestRestoreNodesDropsRemovedNodes checks the restore is a swap, not a merge.
func TestRestoreNodesDropsRemovedNodes(t *testing.T) {
	grid, _ := newRestoreGrid(t)
	addBang(grid, 1, 1)
	doc := grid.Serialize()

	addBang(grid, 2, 2)
	grid.RestoreNodes(doc)

	if grid.Node(2, 2) != nil {
		t.Error("node added after the snapshot survived the restore")
	}
	if grid.Node(1, 1) == nil {
		t.Error("node in the snapshot is missing after the restore")
	}
}
