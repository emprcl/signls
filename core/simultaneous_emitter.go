package core

type SimultaneousEmitter struct {
	note      Note
	direction Direction
	pulse     uint64
	armed     bool
	triggered bool
}

func (e *SimultaneousEmitter) Activated() bool {
	return e.armed || e.triggered
}

func (e *SimultaneousEmitter) Note() Note {
	return e.note
}

func (e *SimultaneousEmitter) Arm() {
	e.armed = true
}

func (e *SimultaneousEmitter) Trig(g *Grid, x, y int) {
	if e.armed {
		g.Trig(x, y)
		e.triggered = true
		e.armed = false
		e.pulse = g.Pulse
	}
	// TODO: handle length and note off
}

func (e *SimultaneousEmitter) Emit(g *Grid, x, y int) {
	if !e.updated(g.Pulse) && e.triggered {
		g.Emit(x, y, e.direction)
		e.triggered = false
		e.pulse = g.Pulse
	}
}

func (e *SimultaneousEmitter) Direction() Direction {
	return e.direction
}

func (e *SimultaneousEmitter) Symbol() string {
	return "S"
}

func (e *SimultaneousEmitter) Color() string {
	return "165"
}

func (e *SimultaneousEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
