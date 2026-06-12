package music

import (
	"math"
	"math/rand"
	"time"

	"signls/core/theory"
	"signls/midi"
)

const (
	maxKey theory.Key = 127
	minKey theory.Key = 21
)

type KeyValue struct {
	rand *rand.Rand

	key      theory.Key
	nextKey  theory.Key
	lastKey  theory.Key
	interval int
	amount   int

	silent bool
}

func NewKeyValue(key theory.Key) *KeyValue {
	source := rand.NewSource(time.Now().UnixNano())
	return &KeyValue{
		key:  key,
		rand: rand.New(source),
	}
}

func (p *KeyValue) BaseValue() theory.Key {
	return p.key
}

func (p *KeyValue) Value() theory.Key {
	if p.nextKey > 0 {
		return p.nextKey
	}
	return p.key
}

func (p *KeyValue) Last() theory.Key {
	return p.lastKey
}

func (p *KeyValue) Display() string {
	return midi.Note(uint8(p.Value()))
}

func (p *KeyValue) Computed(root theory.Key, scale theory.Scale) theory.Key {
	if p.nextKey > 0 {
		p.key = p.nextKey
		p.nextKey = 0
	}
	if p.amount == 0 {
		p.lastKey = p.key
		return p.lastKey
	}
	key := theory.Key(p.rand.Intn(int(math.Abs(float64(p.amount))) + 1))
	if p.amount > 0 {
		key = p.key + key
	} else {
		key = p.key - key
	}
	interval := key.AllSemitonesFrom(root)
	p.lastKey = p.key.Transpose(root, scale, interval)
	return p.lastKey
}

func (p *KeyValue) SetNext(key theory.Key, root theory.Key) {
	if key < minKey || key > maxKey {
		return
	}
	p.nextKey = key
	if int(key)+p.amount < int(maxKey) {
		p.amount++
	}
	if int(key)+p.amount > int(minKey) {
		p.amount--
	}
	p.interval = p.nextKey.AllSemitonesFrom(root)
}

func (p *KeyValue) Set(key theory.Key) {
	p.key = key
	p.nextKey = 0
}

func (p *KeyValue) RandomAmount() int {
	return p.amount
}

func (p *KeyValue) SetRandomAmount(amount int) {
	if int(p.key)+amount < int(minKey) || int(p.key)+amount > int(maxKey) {
		return
	}
	p.amount = amount
}

func (p *KeyValue) IsSilent() bool {
	return p.silent
}

func (p *KeyValue) SetSilent(silent bool) {
	p.silent = silent
}
