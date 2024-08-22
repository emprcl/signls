package core

type EmitterBehavior interface {
	ArmedOnStart() bool
	EmitDirections(dir Direction, pulse uint64) Direction
	Symbol(dir Direction) string
	Name() string
	Color() string
}

type BaseEmitter struct {
	behavior EmitterBehavior

	direction Direction
	note      *Note

	pulse     uint64
	armed     bool
	triggered bool
	muted     bool
}

func (e *BaseEmitter) Copy() Node {
	newNote := *e.note
	return &BaseEmitter{
		behavior:  e.behavior,
		direction: e.direction,
		armed:     e.armed,
		note:      &newNote,
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
	directions := e.behavior.EmitDirections(e.direction, g.pulse)
	for _, dir := range directions.Decompose() {
		g.Emit(x, y, dir)
	}
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
	return e.behavior.Symbol(e.direction)
}

func (e *BaseEmitter) Name() string {
	return e.behavior.Name()
}

func (e *BaseEmitter) Color() string {
	return e.behavior.Color()
}

func (e *BaseEmitter) Reset() {
	e.pulse = 0
	e.armed = e.behavior.ArmedOnStart()
	e.triggered = false
	e.Note().Stop()
}

func (e *BaseEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
