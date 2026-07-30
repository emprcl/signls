package ui

import (
	"testing"
	"time"

	"signls/core/theory"
	"signls/filesystem"
)

// docWithWidth builds a minimal, distinguishable document. Width is used as an
// easy-to-assert identity for a state: unlike tempo it is part of the document
// proper, so sameDoc doesn't ignore it.
func docWithWidth(width int) filesystem.Grid {
	g := filesystem.NewGrid()
	g.Width = width
	return g
}

// docAt returns a state carrying the document identified by width.
func docAt(width int) state {
	return state{doc: docWithWidth(width)}
}

// newTestHistory returns a history whose clock the test drives, plus the advance
// function. Coalescing is time-based, so it can't be tested against a real clock.
func newTestHistory(width int) (*history, func(time.Duration)) {
	clock := time.Unix(0, 0)
	h := newHistory(docWithWidth(width))
	h.now = func() time.Time { return clock }
	return h, func(d time.Duration) { clock = clock.Add(d) }
}

func TestHistoryUndoRedo(t *testing.T) {
	h := newHistory(docWithWidth(100))
	h.push(docAt(110))
	h.push(docAt(120))

	s, ok := h.undo()
	if !ok || s.doc.Width != 110 {
		t.Fatalf("undo 1: got (%v, %v), want (110, true)", s.doc.Width, ok)
	}
	s, ok = h.undo()
	if !ok || s.doc.Width != 100 {
		t.Fatalf("undo 2: got (%v, %v), want (100, true)", s.doc.Width, ok)
	}
	if _, ok := h.undo(); ok {
		t.Fatalf("undo past oldest should fail")
	}

	s, ok = h.redo()
	if !ok || s.doc.Width != 110 {
		t.Fatalf("redo 1: got (%v, %v), want (110, true)", s.doc.Width, ok)
	}
	s, ok = h.redo()
	if !ok || s.doc.Width != 120 {
		t.Fatalf("redo 2: got (%v, %v), want (120, true)", s.doc.Width, ok)
	}
	if _, ok := h.redo(); ok {
		t.Fatalf("redo past newest should fail")
	}
}

// TestHistoryPushDiscardsRedoBranch covers the case: undo several times, edit,
// then redo. The redo branch must be discarded and redo must be a no-op.
func TestHistoryPushDiscardsRedoBranch(t *testing.T) {
	h := newHistory(docWithWidth(100))
	h.push(docAt(110)) // A
	h.push(docAt(120)) // B

	h.undo() // back to A (110)
	h.undo() // back to baseline (100)

	h.push(docAt(130)) // new edit C from baseline

	// A and B must be gone; redo has nothing to replay.
	if _, ok := h.redo(); ok {
		t.Fatalf("redo after divergent edit should be a no-op")
	}

	// Undo must now walk back to the baseline, not the discarded branch.
	s, ok := h.undo()
	if !ok || s.doc.Width != 100 {
		t.Fatalf("undo after divergent edit: got (%v, %v), want (100, true)", s.doc.Width, ok)
	}
}

// TestHistoryPushIgnoresIdenticalState ensures no-op edits don't create undo
// steps.
func TestHistoryPushIgnoresIdenticalState(t *testing.T) {
	h := newHistory(docWithWidth(100))
	h.push(docAt(100)) // identical to baseline

	if _, ok := h.undo(); ok {
		t.Fatalf("identical push should not create an undo step")
	}
}

// TestHistoryPushIgnoresSequencerMusicalChanges is the meta command case: the
// clock goroutine changed the root key, scale and tempo behind the user's back,
// so a commit whose document is otherwise unchanged is still a no-op edit. Were
// it recorded, undo would revert a change the user never made.
func TestHistoryPushIgnoresSequencerMusicalChanges(t *testing.T) {
	h := newHistory(docWithWidth(100))

	s := docAt(100)
	s.doc.Key, s.doc.Scale, s.doc.Tempo = 72, 4, 180
	h.push(s)

	if _, ok := h.undo(); ok {
		t.Fatalf("a musical change the user didn't make should not create an undo step")
	}
}

// TestHistoryPushRecordsUserMusicalEdits is the other half: the same document
// with an explicit musicalEdit is a real edit, and undo/redo replay the recorded
// values rather than the snapshot's.
func TestHistoryPushRecordsUserMusicalEdits(t *testing.T) {
	h := newHistory(docWithWidth(100))

	s := docAt(100)
	s.edit = musicalEdit{
		fields: touchKey,
		before: musical{key: theory.Key(60)},
		after:  musical{key: theory.Key(61)},
	}
	h.push(s)

	undone, ok := h.undo()
	if !ok {
		t.Fatalf("a user musical edit should create an undo step")
	}
	if undone.fields != touchKey || undone.musical.key != theory.Key(60) {
		t.Fatalf("undo: got (%v, %v), want (touchKey, 60)", undone.fields, undone.musical.key)
	}

	redone, ok := h.redo()
	if !ok {
		t.Fatalf("redo should replay the musical edit")
	}
	if redone.fields != touchKey || redone.musical.key != theory.Key(61) {
		t.Fatalf("redo: got (%v, %v), want (touchKey, 61)", redone.fields, redone.musical.key)
	}
}

// TestHistoryCoalescesGesture covers a held key: a burst of commits sharing a
// gesture is one undo step, not one per repeat.
func TestHistoryCoalescesGesture(t *testing.T) {
	h, advance := newTestHistory(0)

	for i := 1; i <= 20; i++ {
		s := docAt(i)
		s.gesture = "param:velocity"
		h.push(s)
		advance(10 * time.Millisecond)
	}

	if len(h.states) != 2 {
		t.Fatalf("states: got %d, want 2 (baseline + one coalesced step)", len(h.states))
	}
	s, ok := h.undo()
	if !ok || s.doc.Width != 0 {
		t.Fatalf("undo: got (%v, %v), want (0, true) — one undo should span the burst", s.doc.Width, ok)
	}
}

// TestHistoryCoalesceKeepsFirstBefore checks a coalesced musical burst undoes to
// where the burst started, not to its penultimate value.
func TestHistoryCoalesceKeepsFirstBefore(t *testing.T) {
	h, advance := newTestHistory(0)

	for i := 1; i <= 5; i++ {
		s := docAt(0)
		s.gesture = gestureTempo
		s.edit = musicalEdit{
			fields: touchTempo,
			before: musical{tempo: float64(119 + i)},
			after:  musical{tempo: float64(120 + i)},
		}
		h.push(s)
		advance(10 * time.Millisecond)
	}

	s, ok := h.undo()
	if !ok || s.musical.tempo != 120 {
		t.Fatalf("undo: got (%v, %v), want (120, true)", s.musical.tempo, ok)
	}
}

// TestHistoryCoalesceDropsRoundTrip covers a gesture that ends where it started —
// toggling a node direction twice, walking a value up then back down. The step
// would undo to nothing, so it shouldn't exist.
func TestHistoryCoalesceDropsRoundTrip(t *testing.T) {
	h, advance := newTestHistory(20)

	away := docAt(30)
	away.gesture = "direction:1:1:1:1"
	h.push(away)
	advance(10 * time.Millisecond)

	back := docAt(20)
	back.gesture = "direction:1:1:1:1"
	h.push(back)

	if len(h.states) != 1 {
		t.Fatalf("states: got %d, want 1 — a round trip is not an edit", len(h.states))
	}
	if _, ok := h.undo(); ok {
		t.Fatalf("a round trip should leave nothing to undo")
	}
}

// TestHistoryGestureExpires checks a gesture closes once the user goes quiet, so
// two deliberate edits of the same parameter stay two undo steps.
func TestHistoryGestureExpires(t *testing.T) {
	h, advance := newTestHistory(0)

	first := docAt(1)
	first.gesture = "param:velocity"
	h.push(first)

	advance(coalesceWindow)

	second := docAt(2)
	second.gesture = "param:velocity"
	h.push(second)

	if len(h.states) != 3 {
		t.Fatalf("states: got %d, want 3 (baseline + two steps)", len(h.states))
	}
}

// TestHistoryGestureDoesNotOverwriteRedoBranch checks that after an undo, a
// commit sharing the gesture of the state we landed on appends instead of
// overwriting it — otherwise a stale redo branch would survive the divergence.
func TestHistoryGestureDoesNotOverwriteRedoBranch(t *testing.T) {
	h, _ := newTestHistory(0)

	first := docAt(1)
	first.gesture = "param:velocity"
	h.push(first)

	advanced := docAt(2)
	h.push(advanced)

	h.undo() // back onto first, whose gesture is still open in time

	next := docAt(3)
	next.gesture = "param:velocity"
	h.push(next)

	if _, ok := h.redo(); ok {
		t.Fatalf("redo branch survived a divergent edit")
	}
	s, ok := h.undo()
	if !ok || s.doc.Width != 1 {
		t.Fatalf("undo: got (%v, %v), want (1, true)", s.doc.Width, ok)
	}
}

// TestHistoryRestoresCursor checks a step carries the position of the change, so
// undo puts the cursor on what it just reverted.
func TestHistoryRestoresCursor(t *testing.T) {
	h := newHistory(docWithWidth(100))

	s := docAt(110)
	s.cursorX, s.cursorY = 4, 7
	s.selectionX, s.selectionY = 5, 9
	h.push(s)

	undone, _ := h.undo()
	if undone.cursorX != 4 || undone.cursorY != 7 || undone.selectionX != 5 || undone.selectionY != 9 {
		t.Fatalf("undo cursor: got (%d,%d)-(%d,%d), want (4,7)-(5,9)",
			undone.cursorX, undone.cursorY, undone.selectionX, undone.selectionY)
	}
}

// TestHistoryRebase checks a non-edit change to the live document updates the
// current state in place rather than creating a step.
func TestHistoryRebase(t *testing.T) {
	h := newHistory(docWithWidth(20))
	h.push(docAt(30))

	h.rebase(docWithWidth(40))

	if len(h.states) != 2 {
		t.Fatalf("states: got %d, want 2 — rebase must not create a step", len(h.states))
	}
	if got := h.states[h.index].doc.Width; got != 40 {
		t.Fatalf("current doc width: got %d, want 40", got)
	}
	// The next commit must diff against the rebased document, not the old one.
	h.push(docAt(40))
	if len(h.states) != 2 {
		t.Fatalf("states: got %d, want 2 — a commit matching the rebased doc is a no-op", len(h.states))
	}
}

// TestHistoryLimit ensures the oldest states are dropped past the limit while
// the newest state stays current.
func TestHistoryLimit(t *testing.T) {
	h := newHistory(docWithWidth(0))
	for i := 1; i <= historyLimit+50; i++ {
		h.push(docAt(i))
	}
	if len(h.states) != historyLimit {
		t.Fatalf("states length: got %d, want %d", len(h.states), historyLimit)
	}
	if h.states[h.index].doc.Width != historyLimit+50 {
		t.Fatalf("current state: got %v, want %v", h.states[h.index].doc.Width, historyLimit+50)
	}
}
