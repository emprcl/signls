package sequencer

import "cykl/midi"

const (
	defaultNote     uint8 = 60
	defaultVelocity uint8 = 100
)

type Track struct {
	midi     midi.Midi
	trig     chan struct{}
	tick     chan struct{}
	done     chan struct{}
	chord    []uint8
	pulse    int
	Steps    int
	length   int
	velocity uint8
	channel  uint8
}

func NewTrack(midi midi.Midi) *Track {
	t := &Track{
		midi:     midi,
		Steps:    16,
		chord:    []uint8{defaultNote},
		length:   pulsesPerStep,
		velocity: defaultVelocity,
	}

	t.trig = make(chan struct{})
	t.tick = make(chan struct{})
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
	return t.pulse / pulsesPerStep
}

func (t *Track) Tick() {
	t.trig <- struct{}{}
}

func (t *Track) Reset() {
	t.pulse = 0
}

func (t *Track) trigger() {
	go func(track *Track) {
		start := 0
		for _, note := range t.chord {
			track.midi.NoteOn(0, t.channel, note, t.velocity)
		}
		for {
			if start >= t.length {
				for _, note := range t.chord {
					track.midi.NoteOff(0, t.channel, note)
				}
				break
			}
			<-track.tick
			start++
		}
	}(t)
	t.pulse = (t.pulse + 1) % (pulsesPerStep * t.Steps)
}
