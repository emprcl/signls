package node

import (
	"signls/core/common"
)

type Signal struct {
	direction common.Direction
	pulse     uint64
	ticks     uint64
	division  int
}

func NewSignal(direction common.Direction, pulse uint64, division int) *Signal {
	return &Signal{
		direction: direction,
		pulse:     pulse,
		division:  division,
	}
}

func (s *Signal) MustMove(pulse uint64) bool {
	if !s.updated(pulse) && s.ticks >= uint64(s.division) {
		s.pulse = pulse
		s.ticks = 0
		return true
	}
	return false
}

func (s *Signal) Direction() common.Direction {
	return s.direction
}

func (s *Signal) Division() int {
	return s.division
}

func (s *Signal) Tick() {
	s.ticks++
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

func (s *Signal) Reset() {
	s.pulse = 0
	s.ticks = 0
}

func (s *Signal) updated(pulse uint64) bool {
	return s.pulse == pulse
}
