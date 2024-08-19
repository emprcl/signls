package core

type BaseEmitter struct {
	note      *Note
	direction Direction
	pulse     uint64
	armed     bool
	triggered bool
	muted     bool

	EmitterBehavior
}

func (e *BaseEmitter) Copy() Node {
	newNote := *e.note
	return &BaseEmitter{
		direction:       e.direction,
		armed:           e.armed,
		note:            &newNote,
		EmitterBehavior: e.EmitterBehavior,
	}
}

func (e *BaseEmitter) Activated() bool {
	return e.armed || e.triggered
}

func (e *BaseEmitter) Note() *Note {
	return e.note
}

func (e *BaseEmitter) Arm() {
	e.armed = true
}

func (e *BaseEmitter) SetMute(mute bool) {
	e.note.Stop()
	e.muted = mute
}

func (e *BaseEmitter) Muted() bool {
	return e.muted
}

func (e *BaseEmitter) Trig(key Key, scale Scale, pulse uint64) {
	if !e.updated(pulse) {
		e.note.Tick()
	}
	if !e.armed {
		return
	}
	if !e.muted {
		e.note.Play(key, scale)
	}
	e.triggered = true
	e.armed = false
	e.pulse = pulse
}

func (e *BaseEmitter) Emit(g *Grid, x, y int) {
	if e.updated(g.pulse) || !e.triggered {
		return
	}
	e.EmitterBehavior.Emit(g, e.direction, x, y)
	e.triggered = false
	e.pulse = g.pulse
}

func (e *BaseEmitter) Direction() Direction {
	return e.direction
}

func (e *BaseEmitter) SetDirection(dir Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

func (e *BaseEmitter) Symbol() string {
	return e.EmitterBehavior.Symbol(e.direction)
}

func (e *BaseEmitter) Name() string {
	return e.EmitterBehavior.Name()
}

func (e *BaseEmitter) Color() string {
	return e.EmitterBehavior.Color()
}

func (e *BaseEmitter) Reset() {
	e.pulse = 0
	e.armed = e.EmitterBehavior.ArmedOnStart()
	e.triggered = false
	e.Note().Stop()
}

func (e *BaseEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
