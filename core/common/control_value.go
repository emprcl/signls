package common

import (
	"math"
	"math/rand"
	"time"
)

type Number interface {
	int | uint8
}

type ControlValue[T Number] struct {
	val      T
	last     T
	min, max T
	amount   int
	rand     *rand.Rand
}

func NewControlValue[T Number](value T, min T, max T) *ControlValue[T] {
	source := rand.NewSource(time.Now().UnixNano())
	return &ControlValue[T]{
		val:  value,
		last: value,
		min:  min,
		max:  max,
		rand: rand.New(source),
	}
}

func (p *ControlValue[T]) Value() T {
	return p.val
}

func (p *ControlValue[T]) Computed() T {
	if p.amount == 0 {
		p.last = p.val
		return p.last
	}
	value := T(p.rand.Intn(int(math.Abs(float64(p.amount))) + 1))
	if p.amount > 0 {
		value = p.val + value
	} else {
		value = p.val - value
	}
	p.last = max(min(value, p.max), p.min)
	return p.last
}

func (p *ControlValue[T]) Last() T {
	return p.last
}

func (p *ControlValue[T]) Set(value T) {
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

func (p *ControlValue[T]) RandomAmount() int {
	return p.amount
}

func (p *ControlValue[T]) SetRandomAmount(amount int) {
	if int(p.val)+amount < int(p.min) || int(p.val)+amount > int(p.max) {
		return
	}
	p.amount = amount
}

func (p *ControlValue[T]) SetMin(min T) {
	if p.val < min {
		p.val = min
	}
	p.min = min
}

func (p *ControlValue[T]) SetMax(max T) {
	if p.val > max {
		p.val = max
	}
	p.max = max
}
