package node

import (
	"cykl/core/common"
)

type TeleportEmitter struct {
	activated    int
	destinationX int
	destinationY int
}

func NewTeleportEmitter(direction common.Direction, x, y int) *TeleportEmitter {
	return &TeleportEmitter{
		destinationX: x,
		destinationY: y,
	}
}

func (e *TeleportEmitter) Copy(dx, dy int) common.Node {
	return &TeleportEmitter{
		destinationX: e.destinationX + dx,
		destinationY: e.destinationY + dy,
	}
}

func (e *TeleportEmitter) Activated() bool {
	return e.activated > 0
}

func (s *TeleportEmitter) Direction() common.Direction {
	return common.NONE
}

func (s *TeleportEmitter) SetDirection(dir common.Direction) {}

func (e *TeleportEmitter) Teleport() (int, int) {
	e.activated = common.PulsesPerStep + 1
	return e.Destination()
}

func (e *TeleportEmitter) Destination() (int, int) {
	return e.destinationX, e.destinationY
}

func (e *TeleportEmitter) SetDestination(x, y int) {
	e.destinationX, e.destinationY = x, y
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
	return "124"
}
