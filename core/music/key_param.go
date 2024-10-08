package music

import (
	"cykl/midi"
	"math"
	"math/rand"
	"time"
)

const (
	maxKey Key = 127
	minKey Key = 21
)

type KeyParam struct {
	rand *rand.Rand

	key      Key
	nextKey  Key
	lastKey  Key
	interval int
	amount   int

	silent bool
}

func NewKeyParam(key Key) *KeyParam {
	source := rand.NewSource(time.Now().UnixNano())
	return &KeyParam{
		key:  key,
		rand: rand.New(source),
	}
}

func (p *KeyParam) Value() Key {
	if p.nextKey > 0 {
		return p.nextKey
	}
	return p.key
}

func (p *KeyParam) Last() Key {
	return p.lastKey
}

func (p *KeyParam) Display() string {
	return midi.Note(uint8(p.Value()))
}

func (p *KeyParam) Computed(root Key, scale Scale) Key {
	if p.nextKey > 0 {
		p.key = p.nextKey
		p.nextKey = 0
	}
	if p.amount == 0 {
		return p.key
	}
	key := Key(p.rand.Intn(int(math.Abs(float64(p.amount))) + 1))
	if p.amount > 0 {
		key = p.key + Key(key)
	} else {
		key = p.key - Key(key)
	}
	interval := key.AllSemitonesFrom(root)
	p.lastKey = p.key.Transpose(root, scale, interval)
	return p.lastKey
}

func (p *KeyParam) SetNext(key Key, root Key) {
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

func (p *KeyParam) Set(key Key) {
	p.key = key
	p.nextKey = 0
}

func (p *KeyParam) Transpose(root Key, scale Scale) {
	p.SetNext(p.key.Transpose(root, scale, p.interval), root)
}

func (p *KeyParam) RandomAmount() int {
	return p.amount
}

func (p *KeyParam) SetRandomAmount(amount int) {
	if int(p.key)+amount < int(minKey) || int(p.key)+amount > int(maxKey) {
		return
	}
	p.amount = amount
}

func (p *KeyParam) IsChanging() bool {
	return p.nextKey > 0
}

func (p *KeyParam) IsSilent() bool {
	return p.silent
}

func (p *KeyParam) SetSilent(silent bool) {
	p.silent = silent
}

func (p *KeyParam) Name() string {
	return "key"
}

func (p *KeyParam) Symbol() string {
	if p.silent {
		return "\u0353"
	}
	if p.amount != 0 {
		return "\u033c"
	}
	return ""
}
