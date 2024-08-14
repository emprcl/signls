package core

type SimpleEmitter struct {
	note      Note
	direction Direction
	pulse     uint64
	armed     bool
	triggered bool
	muted     bool
}

func NewSimpleEmitter(direction Direction) *SimpleEmitter {
	return &SimpleEmitter{
		direction: direction,
		note: Note{
			Channel:  uint8(0),
			Key:      60,
			Velocity: 100,
		},
	}
}

func (e *SimpleEmitter) Copy() Node {
	return &SimpleEmitter{
		direction: e.direction,
		note:      e.note,
	}
}

func (e *SimpleEmitter) Activated() bool {
	return e.armed || e.triggered
}

func (e *SimpleEmitter) Note() Note {
	return e.note
}

func (e *SimpleEmitter) Arm() {
	e.armed = true
}

func (e *SimpleEmitter) SetMute(mute bool) {
	e.muted = mute
}

func (e *SimpleEmitter) Muted() bool {
	return e.muted
}

func (e *SimpleEmitter) Trig(g *Grid, x, y int) {
	if !e.armed {
		return
	}
	if !e.muted {
		g.Trig(x, y)
	}
	e.triggered = true
	e.armed = false
	e.pulse = g.Pulse
	// TODO: handle length and note off
}

func (e *SimpleEmitter) Emit(g *Grid, x, y int) {
	if e.updated(g.Pulse) || !e.triggered {
		return
	}
	for _, dir := range e.direction.Decompose() {
		g.Emit(x, y, dir)
	}
	e.triggered = false
	e.pulse = g.Pulse
}

func (e *SimpleEmitter) Direction() Direction {
	return e.direction
}

func (e *SimpleEmitter) SetDirection(dir Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

func (e *SimpleEmitter) Symbol() string {
	return "S"
}

func (e *SimpleEmitter) Name() string {
	return "E Simple"
}

func (e *SimpleEmitter) Color() string {
	return "177"
}

func (e *SimpleEmitter) Reset() {
	e.pulse = 0
	e.armed = false
	e.triggered = false
}

func (e *SimpleEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
