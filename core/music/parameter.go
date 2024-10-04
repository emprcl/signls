package music

import "math/rand"

type KeyParam struct {
	key    Key
	rand   *rand.Rand
	amount int
}

func (p *KeyParam) Value() Key {
	return p.key
}

func (p *KeyParam) Set(value Key) {
	p.key = value
}

func (p *KeyParam) RandomAmount() int {
	return p.amount
}

func (p *KeyParam) SetRandomAmount(amount int) {
	p.amount = amount
}

type MidiParam struct {
	val    uint8
	rand   *rand.Rand
	amount int
}

func (p *MidiParam) Value() uint8 {
	return p.val
}

func (p *MidiParam) Set(value uint8) {
	p.val = value
}

func (p *MidiParam) RandomAmount() int {
	return p.amount
}

func (p *MidiParam) SetRandomAmount(amount int) {
	p.amount = amount
}
