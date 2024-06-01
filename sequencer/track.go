package sequencer

import "cykl/midi"

type Track struct {
	midi midi.Midi

	// Each track starts a goroutine to handle its pulse progression and step
	// triggering, by using the trig chan at each clock tick.
	// On track removal, we use the done chan to terminate the goroutine.
	trig chan struct{}
	done chan struct{}

	// The pulse defines the current position of the playhead in the track.
	// Each time the clock ticks, we increment the pulse.
	// pulse ranges from 0 to len(steps) * pulsesPerStep (check clock.go).
	// Because each track can have a different number of steps, track pulses
	// are not always synchronized.
	pulse int

	Steps int
}

func NewTrack(midi midi.Midi) *Track {
	t := &Track{
		midi:  midi,
		Steps: 16,
	}

	t.trig = make(chan struct{})
	t.done = make(chan struct{})
	go func(track *Track) {
		for {
			select {
			case <-track.trig:
				track.trigger()
			case <-track.done:
				return
			}
		}
	}(t)

	return t
}

func (t *Track) CurrentStep() int {
	return t.pulse / 6
}

func (t *Track) Tick() {
	t.trig <- struct{}{}
}

func (t *Track) Reset() {
	t.pulse = 0
}

func (t *Track) trigger() {
	t.pulse = (t.pulse + 1) % (6 * t.Steps)
}
