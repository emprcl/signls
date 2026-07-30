package ui

import (
	"sync"
	"time"
)

// saveDebounce is how long the saver waits after the last change before writing
// to disk. It coalesces bursts of edits — such as a held key repeating, which
// would otherwise rewrite the entire bank file on every repeat — into a single
// write once the edits go quiet.
const saveDebounce = 400 * time.Millisecond

// saver coalesces frequent save requests into a single debounced disk write.
// Each request (re)starts the debounce window; the save runs once requests stop
// for saveDebounce. The actual save is delegated to a caller-provided function,
// which must be safe to call from a background goroutine (the timer fires on
// its own goroutine).
type saver struct {
	mu       sync.Mutex
	timer    *time.Timer
	interval time.Duration
	save     func()
}

func newSaver(interval time.Duration, save func()) *saver {
	return &saver{interval: interval, save: save}
}

// request schedules a save, restarting the debounce window. It never blocks.
func (s *saver) request() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer == nil {
		s.timer = time.AfterFunc(s.interval, s.save)
		return
	}
	s.timer.Reset(s.interval)
}

// flush runs a pending save right now instead of waiting out the debounce. Used
// before switching bank slots: the save writes the live grid into whichever slot
// is active when it fires, so a pending one must land before the active slot
// changes or the edit that scheduled it is written to the wrong slot — that is,
// lost.
func (s *saver) flush() {
	s.mu.Lock()
	pending := s.timer != nil && s.timer.Stop()
	s.mu.Unlock()
	if pending {
		s.save()
	}
}

// stop cancels any pending save. Used on quit, where the final save is done
// synchronously instead.
func (s *saver) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
	}
}
