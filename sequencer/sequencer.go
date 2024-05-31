package sequencer

import (
	"cykl/midi"
)

const (
	defaultTempo float64 = 120.0
)

type Sequencer struct {
	midi midi.Midi

	Clock   *clock
	Playing bool
	Pulse   int
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
	if !s.Playing {
		s.reset()
	}
}

func (s *Sequencer) CurrentStep() int {
	return s.Pulse / 6
}

func (s *Sequencer) reset() {
	s.Pulse = 0
}

func (s *Sequencer) tick() {
	if !s.Playing {
		return
	}
	s.Pulse = (s.Pulse + 1) % (16 * 6)
}
