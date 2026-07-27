package ui

import (
	"testing"

	"signls/filesystem"
)

// gridWithTempo builds a minimal, distinguishable serialized grid. Tempo is used
// as an easy-to-assert identity for a state.
func gridWithTempo(tempo float64) filesystem.Grid {
	g := filesystem.NewGrid()
	g.Tempo = tempo
	return g
}

func TestHistoryUndoRedo(t *testing.T) {
	h := newHistory(gridWithTempo(100))
	h.push(gridWithTempo(110))
	h.push(gridWithTempo(120))

	state, ok := h.undo()
	if !ok || state.Tempo != 110 {
		t.Fatalf("undo 1: got (%v, %v), want (110, true)", state.Tempo, ok)
	}
	state, ok = h.undo()
	if !ok || state.Tempo != 100 {
		t.Fatalf("undo 2: got (%v, %v), want (100, true)", state.Tempo, ok)
	}
	if _, ok := h.undo(); ok {
		t.Fatalf("undo past oldest should fail")
	}

	state, ok = h.redo()
	if !ok || state.Tempo != 110 {
		t.Fatalf("redo 1: got (%v, %v), want (110, true)", state.Tempo, ok)
	}
	state, ok = h.redo()
	if !ok || state.Tempo != 120 {
		t.Fatalf("redo 2: got (%v, %v), want (120, true)", state.Tempo, ok)
	}
	if _, ok := h.redo(); ok {
		t.Fatalf("redo past newest should fail")
	}
}

// TestHistoryPushDiscardsRedoBranch covers the case: undo several times, edit,
// then redo. The redo branch must be discarded and redo must be a no-op.
func TestHistoryPushDiscardsRedoBranch(t *testing.T) {
	h := newHistory(gridWithTempo(100))
	h.push(gridWithTempo(110)) // A
	h.push(gridWithTempo(120)) // B

	h.undo() // back to A (110)
	h.undo() // back to baseline (100)

	h.push(gridWithTempo(130)) // new edit C from baseline

	// A and B must be gone; redo has nothing to replay.
	if _, ok := h.redo(); ok {
		t.Fatalf("redo after divergent edit should be a no-op")
	}

	// Undo must now walk back to the baseline, not the discarded branch.
	state, ok := h.undo()
	if !ok || state.Tempo != 100 {
		t.Fatalf("undo after divergent edit: got (%v, %v), want (100, true)", state.Tempo, ok)
	}
}

// TestHistoryPushIgnoresIdenticalState ensures no-op edits don't create undo
// steps.
func TestHistoryPushIgnoresIdenticalState(t *testing.T) {
	h := newHistory(gridWithTempo(100))
	h.push(gridWithTempo(100)) // identical to baseline

	if _, ok := h.undo(); ok {
		t.Fatalf("identical push should not create an undo step")
	}
}

// TestHistoryLimit ensures the oldest states are dropped past the limit while
// the newest state stays current.
func TestHistoryLimit(t *testing.T) {
	h := newHistory(gridWithTempo(0))
	for i := 1; i <= historyLimit+50; i++ {
		h.push(gridWithTempo(float64(i)))
	}
	if len(h.states) != historyLimit {
		t.Fatalf("states length: got %d, want %d", len(h.states), historyLimit)
	}
	if h.states[h.index].Tempo != float64(historyLimit+50) {
		t.Fatalf("current state: got %v, want %v", h.states[h.index].Tempo, historyLimit+50)
	}
}
