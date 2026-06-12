package ui

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestSaverCoalesces verifies that a burst of requests within the debounce
// window results in a single save, and that a later request triggers another.
func TestSaverCoalesces(t *testing.T) {
	var saves atomic.Int64
	const interval = 20 * time.Millisecond
	s := newSaver(interval, func() { saves.Add(1) })

	// A burst of requests spaced well under the interval must coalesce.
	for i := 0; i < 10; i++ {
		s.request()
		time.Sleep(interval / 10)
	}
	time.Sleep(3 * interval)
	if got := saves.Load(); got != 1 {
		t.Fatalf("burst of requests: got %d saves, want 1", got)
	}

	// A new request after the window triggers a second save.
	s.request()
	time.Sleep(3 * interval)
	if got := saves.Load(); got != 2 {
		t.Fatalf("after second request: got %d saves, want 2", got)
	}
}

// TestSaverStop verifies that a pending save is cancelled by stop.
func TestSaverStop(t *testing.T) {
	var saves atomic.Int64
	const interval = 20 * time.Millisecond
	s := newSaver(interval, func() { saves.Add(1) })

	s.request()
	s.stop()
	time.Sleep(3 * interval)
	if got := saves.Load(); got != 0 {
		t.Fatalf("stopped saver: got %d saves, want 0", got)
	}
}
