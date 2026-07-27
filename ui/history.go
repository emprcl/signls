package ui

import (
	"reflect"

	"signls/filesystem"
)

// historyLimit caps how many states the undo history keeps. Each state is a
// serialized grid (a few KB at most), so this is a generous, cheap bound.
const historyLimit = 100

// history is a linear undo/redo timeline of serialized grid states. It stores
// full snapshots (memento style) rather than invertible commands, which is cheap
// here because a whole grid already serializes to a small, deep-copied value.
//
// states holds the timeline oldest-first; index points at the current state.
// Anything after index is a redo branch. Recording a new state discards that
// branch — diverging from a past point abandons the abandoned future.
type history struct {
	states []filesystem.Grid
	index  int
}

// newHistory returns a history seeded with an initial baseline state, so the
// user can always undo back to where they started.
func newHistory(initial filesystem.Grid) *history {
	return &history{
		states: []filesystem.Grid{initial},
		index:  0,
	}
}

// push records a newly committed state. Identical states (no-op edits) are
// ignored so they don't create empty undo steps. Any redo branch ahead of the
// cursor is discarded. When the limit is exceeded the oldest state is dropped.
func (h *history) push(state filesystem.Grid) {
	if reflect.DeepEqual(state, h.states[h.index]) {
		return
	}
	h.states = append(h.states[:h.index+1], state)
	if len(h.states) > historyLimit {
		h.states = h.states[len(h.states)-historyLimit:]
	}
	h.index = len(h.states) - 1
}

// undo moves one state back and returns it. The bool is false when already at
// the oldest state.
func (h *history) undo() (filesystem.Grid, bool) {
	if h.index == 0 {
		return filesystem.Grid{}, false
	}
	h.index--
	return h.states[h.index], true
}

// redo moves one state forward and returns it. The bool is false when already at
// the newest state.
func (h *history) redo() (filesystem.Grid, bool) {
	if h.index >= len(h.states)-1 {
		return filesystem.Grid{}, false
	}
	h.index++
	return h.states[h.index], true
}

// reset reinitializes the timeline around a new baseline, discarding all
// previous history. Used when switching to a different grid, whose edits form
// their own independent timeline.
func (h *history) reset(state filesystem.Grid) {
	h.states = []filesystem.Grid{state}
	h.index = 0
}
