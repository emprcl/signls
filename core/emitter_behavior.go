package core

import (
	"cykl/midi"
	"fmt"
)

type EmitterBehavior interface {
	ArmedOnStart() bool
	Emit(g *Grid, direction Direction, x, y int)
	Symbol(dir Direction) string
	Name() string
	Color() string
}

type BangEmitter struct{}

func NewBangEmitter(midi midi.Midi, direction Direction, armed bool) *BaseEmitter {
	return &BaseEmitter{
		direction:       direction,
		armed:           armed,
		note:            NewNote(midi),
		EmitterBehavior: BangEmitter{},
	}
}

func (e BangEmitter) ArmedOnStart() bool {
	return true
}

func (e BangEmitter) Emit(g *Grid, direction Direction, x, y int) {
	for _, dir := range direction.Decompose() {
		g.Emit(x, y, dir)
	}
}

func (e BangEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "B", dir.Symbol())
}

func (e BangEmitter) Name() string {
	return "bang"
}

func (e BangEmitter) Color() string {
	return "165"
}

type SpreadEmitter struct{}

func NewSpreadEmitter(midi midi.Midi, direction Direction) *BaseEmitter {
	return &BaseEmitter{
		direction:       direction,
		note:            NewNote(midi),
		EmitterBehavior: SpreadEmitter{},
	}
}

func (e SpreadEmitter) Emit(g *Grid, direction Direction, x, y int) {
	for _, dir := range direction.Decompose() {
		g.Emit(x, y, dir)
	}
}

func (e SpreadEmitter) ArmedOnStart() bool {
	return false
}

func (e SpreadEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "S", dir.Symbol())
}

func (e SpreadEmitter) Name() string {
	return "spread"
}

func (e SpreadEmitter) Color() string {
	return "177"
}
