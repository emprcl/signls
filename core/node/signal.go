package node

import (
	"signls/core/common"
)

type Signal struct {
	direction common.Direction
	pulse     uint64
}

func NewSignal(direction common.Direction, pulse uint64) *Signal {
	return &Signal{
		direction: direction,
		pulse:     pulse,
	}
}

func (s *Signal) MustMove(pulse uint64) bool {
	if !s.updated(pulse) {
		s.pulse = pulse
		return true
	}
	return false
}

func (s *Signal) Direction() common.Direction {
	return s.direction
}

func (s *Signal) SetDirection(dir common.Direction) {
	s.direction = dir
}

func (s *Signal) Activated() bool {
	return true
}

func (s *Signal) Symbol() string {
	return "  "
}

func (s *Signal) Name() string {
	return "signal"
}

func (s *Signal) Color() string {
	return "15"
}

func (s *Signal) updated(pulse uint64) bool {
	return s.pulse == pulse
}
