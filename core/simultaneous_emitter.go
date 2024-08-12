package core

type SimultaneousEmitter struct {
	direction Direction
	activated bool
	updated   bool
}

func (e *SimultaneousEmitter) Emit() {
	e.activated = true
	e.updated = true
}

func (e *SimultaneousEmitter) Update(g *Grid, x, y int) {
	if e.activated && !e.updated {
		g.Emit(x, y, e.direction)
		e.updated = true
		e.activated = false
	} else if e.updated {
		e.updated = false
	}
}

func (e *SimultaneousEmitter) Direction() Direction {
	return e.direction
}

func (e *SimultaneousEmitter) Activated() bool {
	return e.activated
}

func (e *SimultaneousEmitter) Symbol() string {
	return "S"
}

func (e *SimultaneousEmitter) Color() string {
	return "165"
}

func (e *SimultaneousEmitter) Reset() {
	e.updated = false
}
