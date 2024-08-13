package core

type SimultaneousEmitter struct {
	note      Note
	direction Direction
	pulse     int
	triggered bool
	armed     bool
}

func (e *SimultaneousEmitter) Arm() {
	e.armed = true
}

func (e *SimultaneousEmitter) Note() Note {
	return e.note
}

func (e *SimultaneousEmitter) Emit(g *Grid, x, y int) {
	if e.armed && e.triggered {
		g.Emit(x, y, e.direction)
		e.armed = false
		e.triggered = false
	}
}

func (e *SimultaneousEmitter) Trig(g *Grid, x, y int) {
	if e.armed && e.pulse == 0 {
		g.Trig(x, y)
		e.triggered = true
		//e.pulse++
	} else if e.triggered && e.pulse > 0 {
		// TODO: handle length and note off
	}
}

func (e *SimultaneousEmitter) Direction() Direction {
	return e.direction
}

func (e *SimultaneousEmitter) Armed() bool {
	return e.armed
}

func (e *SimultaneousEmitter) Symbol() string {
	return "S"
}

func (e *SimultaneousEmitter) Color() string {
	return "165"
}

func (e *SimultaneousEmitter) Reset() {}
