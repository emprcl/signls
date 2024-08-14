package core

import "time"

const (
	pulsesPerStep       int     = 1
	stepsPerQuarterNote int     = 4
	tempoMin            float64 = 1.0
	tempoMax            float64 = 300.0
	updateBufferSize    int     = 128
)

// clock contains a clock state.
// We use the standard time.ticker as the clock, ticking at the
// standard midi 6 pulses per 16th note (one step).
// The update chan is used to pass new tempo values and recreate a new ticker.
//
// Read more: http://midi.teragonaudio.com/tech/midispec/clock.htm
type clock struct {
	ticker       *time.Ticker
	update       chan float64
	tempo        float64
	shouldUpdate bool
}

func (c *clock) setTempo(tempo float64) {
	if tempo > tempoMax || tempo < tempoMin {
		return
	}
	c.update <- tempo
}

func newClock(tempo float64, tick func()) *clock {
	c := &clock{
		ticker: time.NewTicker(newClockInterval(tempo)),
		update: make(chan float64, updateBufferSize),
		tempo:  tempo,
	}
	go func(c *clock) {
		for {
			select {
			case <-c.ticker.C:
				tick()
				if c.shouldUpdate {
					c.ticker.Reset(newClockInterval(c.tempo))
					c.shouldUpdate = false
				}
			case newTempo := <-c.update:
				c.shouldUpdate = true
				c.tempo = newTempo
			}
		}
	}(c)
	return c
}

func newClockInterval(tempo float64) time.Duration {
	// midi clock: http://midi.teragonaudio.com/tech/midispec/clock.htm
	return time.Duration(1000000*60/(tempo*float64(pulsesPerStep*stepsPerQuarterNote))) * time.Microsecond
}
