package core

type Signal struct {
	direction Direction
	updated   bool
}

func (s *Signal) Move(g *Grid, x, y int) {
	if s.updated {
		s.updated = false
	} else {
		g.Move(x, y, s.direction)
		s.updated = true
	}
}

func (s *Signal) Direction() Direction {
	return s.direction
}

func (s *Signal) Armed() bool {
	return true
}

func (s *Signal) Reset() {
	s.updated = false
}

func (s *Signal) Symbol() string {
	return " "
}

func (s *Signal) Color() string {
	return "15"
}
