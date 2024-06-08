package sequencer

import (
	"cykl/midi"
)

const (
	defaultTempo float64 = 120.0
)

type Sequencer struct {
	midi    midi.Midi
	Clock   *clock
	Tracks  []*Track
	Playing bool
}

// New creates a new sequencer and starts the clock.
func New(midi midi.Midi) *Sequencer {
	seq := &Sequencer{
		midi:    midi,
		Playing: false,
		Tracks: []*Track{
			NewTrack(midi, 0),
			NewTrack(midi, 1),
		},
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

func (s *Sequencer) reset() {
	for _, track := range s.Tracks {
		track.Reset()
	}
}

func (s *Sequencer) tick() {
	if !s.Playing {
		return
	}
	for _, track := range s.Tracks {
		track.Tick()
	}
}
