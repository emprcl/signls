package node

import (
	"cykl/core/common"
)

type TeleportEmitter struct {
	direction    common.Direction
	pulse        uint64
	destinationX int
	destinationY int
}

func NewTeleportEmitter(direction common.Direction, pulse uint64, x, y int) *TeleportEmitter {
	return &TeleportEmitter{
		direction:    direction,
		pulse:        pulse,
		destinationX: x,
		destinationY: y,
	}
}

func (s *TeleportEmitter) Direction() common.Direction {
	return s.direction
}

func (s *TeleportEmitter) SetDirection(dir common.Direction) {
	s.direction = dir
}

func (e *TeleportEmitter) ArmedOnStart() bool {
	return false
}

func (s *TeleportEmitter) Symbol() string {
	return "  "
}

func (s *TeleportEmitter) Name() string {
	return "telep"
}

func (s *TeleportEmitter) Color() string {
	return "15"
}

func (s *TeleportEmitter) updated(pulse uint64) bool {
	return s.pulse == pulse
}
