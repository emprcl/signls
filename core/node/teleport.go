package node

import (
	"cykl/core/common"
)

type TeleportEmitter struct {
	direction    common.Direction
	activated    int
	destinationX int
	destinationY int
}

func NewTeleportEmitter(direction common.Direction, x, y int) *TeleportEmitter {
	return &TeleportEmitter{
		direction:    direction,
		destinationX: x,
		destinationY: y,
	}
}

func (e *TeleportEmitter) Activated() bool {
	return e.activated > 0
}

func (s *TeleportEmitter) Direction() common.Direction {
	return s.direction
}

func (s *TeleportEmitter) SetDirection(dir common.Direction) {
	s.direction = dir
}

func (e *TeleportEmitter) TeleportPosition() (int, int) {
	e.activated = common.PulsesPerStep
	return e.destinationX, e.destinationY
}

func (e *TeleportEmitter) Tick() {
	if e.activated <= 0 {
		return
	}
	e.activated--
}

func (e *TeleportEmitter) Reset() {
	e.activated = 0
}

func (s *TeleportEmitter) Symbol() string {
	return "T⬢"
}

func (s *TeleportEmitter) Name() string {
	return "telep"
}

func (s *TeleportEmitter) Color() string {
	return "88"
}
