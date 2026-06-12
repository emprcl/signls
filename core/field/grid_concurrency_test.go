package field

import (
	"signls/core/common"
	"signls/core/music"
	"signls/core/node"
	"signls/midi"
	"sync"
	"testing"
)

// TestGridConcurrentAccess exercises the grid the way the running app does: the
// clock goroutine drives Update while the ui goroutine renders (reads node
// state) and mutates the grid in response to input. It is meant to be run with
// the race detector (go test -race) — before the grid was guarded by its mutex
// on every access path, this reproduced data races and could panic with an
// index-out-of-range when Resize shrank the grid mid-Update.
func TestGridConcurrentAccess(t *testing.T) {
	m := &midi.Mock{}
	grid := NewGrid(32, 32, m, "")
	device := m.NewDevice("", "")
	grid.AddNode(node.NewBangEmitter(m, &device, common.DOWN|common.RIGHT, true), 7, 7)
	grid.AddNode(node.NewSpreadEmitter(m, &device, common.DOWN), 11, 7)
	grid.AddNode(node.NewSpreadEmitter(m, &device, common.LEFT), 11, 11)
	grid.TogglePlay()

	const iterations = 2000
	var wg sync.WaitGroup

	// Clock goroutine: advance the grid on every pulse.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			grid.Update()
		}
	}()

	// Render goroutine: take a display snapshot and read live control state
	// under the read lock, like View.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cells := grid.Snapshot()
			for y := range cells {
				for x := range cells[y] {
					_ = cells[y][x].Symbol
				}
			}
			grid.Read(func() {
				_ = grid.Pulse()
				_ = grid.QuarterNote()
			})
		}
	}()

	// Input goroutine: mutate the grid (resize, add/remove, trigger) like the ui.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%2 == 0 {
				grid.Resize(16, 16)
			} else {
				grid.Resize(32, 32)
			}
			grid.AddNodeFromSymbol("s", 1, 1)
			grid.TriggerNode(1, 1)
			grid.RemoveNodes(1, 1, 1, 1)
			grid.ShiftKey(1)
			grid.ShiftScale(1)
		}
	}()

	// Param goroutine: edit a node's parameters the way the ui does — mutating
	// node state under Write, and reading it back under Read (as the parameter
	// list is built). This races with the clock's Trig without the locks.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			grid.Write(func() {
				if a, ok := grid.Node(7, 7).(music.Audible); ok {
					a.Note().SetVelocity(uint8(i % 128))
				}
			})
			grid.Read(func() {
				if a, ok := grid.Node(7, 7).(music.Audible); ok {
					_ = a.Note().Velocity.Value()
				}
			})
		}
	}()

	wg.Wait()
}
