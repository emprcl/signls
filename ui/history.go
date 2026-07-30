package ui

import (
	"reflect"
	"time"

	"signls/core/theory"
	"signls/filesystem"
)

const (
	// historyLimit caps how many states one timeline keeps. Each state is a
	// serialized grid (a few KB at most), so this is a generous, cheap bound.
	historyLimit = 100

	// coalesceWindow is how long a gesture stays open. Consecutive commits
	// carrying the same gesture fold into a single undo step until the user goes
	// quiet for this long, so holding a key to walk a parameter from 0 to 127
	// costs one undo step instead of 127 — which would otherwise fill the whole
	// timeline and evict the baseline. Same idea as the saver's debounce.
	coalesceWindow = 500 * time.Millisecond
)

// Gestures label an ongoing edit so consecutive commits can coalesce. The empty
// gesture never coalesces: discrete actions (adding, removing, pasting) each
// deserve their own undo step.
const (
	gestureNone  = ""
	gestureRoot  = "root"
	gestureScale = "scale"
	gestureTempo = "tempo"
)

// touch is a set of grid-level musical fields.
type touch uint8

const (
	touchKey touch = 1 << iota
	touchScale
	touchTempo
)

// musical is the grid-level musical state: root key, scale and tempo. These
// three are written both by the user and by meta commands running on the clock
// goroutine, which is why they can't take part in a snapshot — restoring one
// would revert whatever the sequencer did since the snapshot was taken. They are
// tracked as explicit before/after values instead. See musicalEdit.
type musical struct {
	key   theory.Key
	scale theory.Scale
	tempo float64
}

// musicalEdit records which musical fields a single user action changed, with
// the values on either side of it. Undo replays before, redo replays after, and
// fields the user didn't touch are left alone whatever the sequencer has done to
// them meanwhile.
type musicalEdit struct {
	fields touch
	before musical
	after  musical
}

// merge folds e onto an earlier edit from the same gesture, keeping that edit's
// before values so a single undo walks back past the whole burst. after needs no
// merging: it is always read whole from the live grid, so e's copy is current
// for every field.
func (e musicalEdit) merge(prev musicalEdit) musicalEdit {
	merged := e
	merged.fields = prev.fields | e.fields
	if prev.fields&touchKey != 0 {
		merged.before.key = prev.before.key
	}
	if prev.fields&touchScale != 0 {
		merged.before.scale = prev.before.scale
	}
	if prev.fields&touchTempo != 0 {
		merged.before.tempo = prev.before.tempo
	}
	return merged
}

// state is one point on a timeline: the document as it stood, the musical change
// that got there, and where the user was at the time.
type state struct {
	doc  filesystem.Grid
	edit musicalEdit

	cursorX, cursorY       int
	selectionX, selectionY int

	gesture string
	at      time.Time
}

// step is what undo and redo hand back: the document to restore, the musical
// fields to replay (if any) with their values, and where to put the cursor so
// the change the user just stepped over is on screen.
type step struct {
	doc     filesystem.Grid
	fields  touch
	musical musical

	cursorX, cursorY       int
	selectionX, selectionY int
}

// history is a linear undo/redo timeline for a single bank slot. It stores full
// document snapshots (memento style) rather than invertible commands, which is
// cheap here because a whole grid already serializes to a small, deep-copied
// value. Each bank slot gets its own history: switching banks — including a
// switch a bank meta command performs mid-playback — must not destroy the user's
// undo stack.
//
// states holds the timeline oldest-first; index points at the current state.
// Anything after index is a redo branch. Recording a new state discards that
// branch — diverging from a past point abandons the abandoned future.
type history struct {
	states []state
	index  int

	// now is injectable so tests can drive the coalescing window.
	now func() time.Time
}

// newHistory returns a history seeded with an initial baseline state, so the
// user can always undo back to where they started.
func newHistory(initial filesystem.Grid) *history {
	return &history{
		states: []state{{doc: initial}},
		now:    time.Now,
	}
}

// push records a newly committed state. A state that changes nothing is ignored
// so no-op edits don't create empty undo steps; a state continuing the current
// gesture replaces it; anything else appends and discards the redo branch ahead
// of the cursor. When the limit is exceeded the oldest state is dropped.
func (h *history) push(s state) {
	current := h.states[h.index]
	if s.edit.fields == 0 && sameDoc(s.doc, current.doc) {
		return
	}

	s.at = h.now()
	if h.continues(current, s) {
		s.edit = s.edit.merge(current.edit)
		h.states[h.index] = s
		// A gesture can end back where it started — toggling a direction twice,
		// walking a value up and back down. That leaves a step whose undo does
		// nothing visible, so drop it instead.
		if s.edit.fields == 0 && sameDoc(s.doc, h.states[h.index-1].doc) {
			h.states = h.states[:h.index]
			h.index--
		}
		return
	}

	h.states = append(h.states[:h.index+1], s)
	if len(h.states) > historyLimit {
		h.states = h.states[len(h.states)-historyLimit:]
	}
	h.index = len(h.states) - 1
}

// continues reports whether s extends the gesture that produced current: same
// non-empty gesture, still inside the coalescing window, and no redo branch in
// the way — diverging from a past point must always create a new step rather
// than overwrite one.
func (h *history) continues(current, s state) bool {
	return s.gesture != gestureNone &&
		s.gesture == current.gesture &&
		h.index == len(h.states)-1 &&
		s.at.Sub(current.at) < coalesceWindow
}

// undo moves one state back and returns the step that gets there. The bool is
// false when already at the oldest state. The musical values are the ones from
// before the undone action, and the cursor is the one recorded on it — that's
// where the change being reverted happened.
func (h *history) undo() (step, bool) {
	if h.index == 0 {
		return step{}, false
	}
	undone := h.states[h.index]
	h.index--
	return h.stepTo(h.states[h.index].doc, undone, undone.edit.before), true
}

// redo moves one state forward and returns the step that gets there. The bool is
// false when already at the newest state.
func (h *history) redo() (step, bool) {
	if h.index >= len(h.states)-1 {
		return step{}, false
	}
	h.index++
	redone := h.states[h.index]
	return h.stepTo(redone.doc, redone, redone.edit.after), true
}

// stepTo assembles a step from the document to restore and the state whose edit
// is being replayed in one direction or the other.
func (h *history) stepTo(doc filesystem.Grid, replayed state, values musical) step {
	return step{
		doc:        doc,
		fields:     replayed.edit.fields,
		musical:    values,
		cursorX:    replayed.cursorX,
		cursorY:    replayed.cursorY,
		selectionX: replayed.selectionX,
		selectionY: replayed.selectionY,
	}
}

// rebase replaces the current state's document without creating a new step. Used
// when something other than a user edit changes the live document — the window
// growing the grid to fill itself — so that states[index] keeps describing the
// grid that is actually loaded. Without it the next commit would diff against a
// state the grid isn't in, and an undo of a resize would immediately be undone
// again by the window.
func (h *history) rebase(doc filesystem.Grid) {
	h.states[h.index].doc = doc
}

// sameDoc reports whether two serialized grids are the same document. The
// grid-level musical fields are ignored: meta commands change them from the
// clock goroutine, so a difference there is not a user edit. They stay in the
// stored document because a bank-slot restore writes it back to disk verbatim.
func sameDoc(a, b filesystem.Grid) bool {
	a.Key, a.Scale, a.Tempo = 0, 0, 0
	b.Key, b.Scale, b.Tempo = 0, 0, 0
	return reflect.DeepEqual(a, b)
}
