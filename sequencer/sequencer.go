package sequencer

import (
	"cykl/midi"
	"fmt"
)

const (
	defaultTempo float64 = 120.0
)

type Sequencer struct {
	midi    midi.Midi
	Clock   *clock
	Playing bool
}

// New creates a new sequencer and starts the clock.
func New(midi midi.Midi) *Sequencer {
	seq := &Sequencer{
		midi:    midi,
		Playing: false,
	}

	seq.Clock = newClock(defaultTempo, func() {
		seq.tick()
	})

	return seq
}

// TogglePlay plays or stops the sequencer.
func (s *Sequencer) TogglePlay() {
	s.Playing = !s.Playing
}

func (s *Sequencer) tick() {
	if !s.Playing {
		return
	}

	fmt.Println("TICK")
}
