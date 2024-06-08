package sequencer

import (
	"cykl/midi"
)

const (
	defaultNote     uint8 = 60
	defaultVelocity uint8 = 100
)

type Chord []uint8

type Track struct {
	midi      midi.Midi
	tick      chan struct{}
	done      chan struct{}
	triggered map[int][]Chord
	chord     Chord
	pulse     int
	Steps     int
	length    int
	velocity  uint8
	channel   uint8
}

func NewTrack(midi midi.Midi, channel uint8) *Track {
	t := &Track{
		midi:     midi,
		Steps:    16,
		chord:    Chord{defaultNote},
		length:   pulsesPerStep,
		velocity: defaultVelocity,
		channel:  channel,
	}
	t.start()
	return t
}

func (t *Track) CurrentStep() int {
	return t.pulse / pulsesPerStep
}

func (t *Track) Tick() {
	t.tick <- struct{}{}
}

func (t *Track) Reset() {
	t.pulse = 0
	for _, chords := range t.triggered {
		for _, chord := range chords {
			for _, note := range chord {
				t.midi.NoteOff(0, t.channel, note)
			}
		}
	}
}

func (t *Track) start() {
	t.tick = make(chan struct{})
	t.done = make(chan struct{})
	t.triggered = make(map[int][]Chord)
	go func(track *Track) {
		for {
			select {
			case <-track.tick:
				track.trigger()
			case <-track.done:
				return
			}
		}
	}(t)
}

func (t *Track) shouldTrigger() bool {
	return (t.pulse % pulsesPerStep) == 0
}

func (t *Track) trigger() {
	if chords, shouldStop := t.triggered[t.pulse]; shouldStop {
		for _, chord := range chords {
			for _, note := range chord {
				t.midi.NoteOff(0, t.channel, note)
			}
		}
		delete(t.triggered, t.pulse)
	}

	if t.shouldTrigger() {
		t.note()
	}

	t.pulse = (t.pulse + 1) % (pulsesPerStep * t.Steps)
}

func (t *Track) note() {
	for _, note := range t.chord {
		t.midi.NoteOn(0, t.channel, note, t.velocity)
	}
	stopPulse := (t.pulse + t.length) % (pulsesPerStep * t.Steps)
	t.triggered[stopPulse] = append(t.triggered[stopPulse], t.chord)
}
