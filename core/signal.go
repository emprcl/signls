package core

type Signal struct {
	direction Direction
	pulse     uint64
}

func (s *Signal) Move(g *Grid, x, y int) {
	if !s.updated(g.Pulse) {
		g.Move(x, y, s.direction)
		s.pulse = g.Pulse
	}
}

func (s *Signal) Direction() Direction {
	return s.direction
}

func (s *Signal) Activated() bool {
	return true
}

func (s *Signal) Symbol() string {
	return " "
}

func (s *Signal) Color() string {
	return "15"
}

func (s Signal) updated(pulse uint64) bool {
	return s.pulse == pulse
}
