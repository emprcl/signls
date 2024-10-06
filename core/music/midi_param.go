package music

import (
	"math"
	"math/rand"
	"time"
)

type MidiParam struct {
	val      uint8
	last     uint8
	min, max uint8
	amount   int
	rand     *rand.Rand
}

func NewMidiParam(value uint8, min uint8, max uint8) *MidiParam {
	source := rand.NewSource(time.Now().UnixNano())
	return &MidiParam{
		val:  value,
		min:  min,
		max:  max,
		rand: rand.New(source),
	}
}

func (p *MidiParam) Value() uint8 {
	return p.val
}

func (p *MidiParam) Computed() uint8 {
	if p.amount == 0 {
		return p.val
	}
	value := uint8(p.rand.Intn(int(math.Abs(float64(p.amount))) + 1))
	if p.amount > 0 {
		value = p.val + value
	} else {
		value = p.val - value
	}
	p.last = max(min(value, p.max), p.min)
	return p.last
}

func (p *MidiParam) Last() uint8 {
	return p.last
}

func (p *MidiParam) Set(value uint8) {
	if value < p.min || value > p.max {
		return
	}
	if int(value)+p.amount < int(p.min) {
		p.amount++
	}
	if int(value)+p.amount > int(p.max) {
		p.amount--
	}
	p.val = value
	p.last = value
}

func (p *MidiParam) RandomAmount() int {
	return p.amount
}

func (p *MidiParam) SetRandomAmount(amount int) {
	if int(p.val)+amount < int(p.min) || int(p.val)+amount > int(p.max) {
		return
	}
	p.amount = amount
}
