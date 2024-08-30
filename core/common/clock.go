package common

import "time"

// Constants defining the timing and buffer parameters
const (
	PulsesPerStep       int = 4
	StepsPerQuarterNote int = 4

	tempoMin         float64 = 1.0
	tempoMax         float64 = 300.0
	updateBufferSize int     = 128
)

// clock manages the timing for MIDI playback, using a standard time.Ticker
// to generate clock pulses. It provides functionality to update the tempo dynamically.
// The update channel is used for receiving new tempo values and adjusting the ticker accordingly.
//
// Read more: http://midi.teragonaudio.com/tech/midispec/clock.htm
type Clock struct {
	ticker       *time.Ticker // Ticker to generate regular clock pulses
	update       chan float64 // Channel for updating the tempo
	tempo        float64      // Current tempo in BPM
	shouldUpdate bool         // Flag to indicate if the ticker should be updated
}

// setTempo updates the tempo of the clock. It ensures the new tempo is within the defined range.
// If the tempo is valid, it sends the new tempo to the update channel.
func (c *Clock) SetTempo(tempo float64) {
	if tempo > tempoMax || tempo < tempoMin {
		return
	}
	c.update <- tempo
}

// Tempo returns the tempo of the clock.
func (c *Clock) Tempo() float64 {
	return c.tempo
}

// NewClock creates and initializes a new clock instance with the specified tempo
// and a callback function that is called on each tick. It starts a goroutine to
// manage the clock ticks and tempo updates.
func NewClock(tempo float64, tick func()) *Clock {
	c := &Clock{
		ticker: time.NewTicker(newClockInterval(tempo)),
		update: make(chan float64, updateBufferSize),
		tempo:  tempo,
	}
	go func(c *Clock) {
		for {
			select {
			case <-c.ticker.C: // On each tick, execute the callback
				tick()
				if c.shouldUpdate {
					c.ticker.Reset(newClockInterval(c.tempo)) // Reset the ticker with the new interval
					c.shouldUpdate = false
				}
			case newTempo := <-c.update: // On receiving a new tempo value
				c.shouldUpdate = true
				c.tempo = newTempo
			}
		}
	}(c)
	return c
}

// newClockInterval calculates the duration of each tick based on the current tempo.
// The duration is computed to ensure that the ticker emits ticks at the correct rate
// for the specified tempo.
func newClockInterval(tempo float64) time.Duration {
	// midi clock: http://midi.teragonaudio.com/tech/midispec/clock.htm
	return time.Duration(1000000*60/(tempo*float64(PulsesPerStep*StepsPerQuarterNote))) * time.Microsecond
}
