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
	key      Key
	nextKey  Key
	lastKey  Key
	interval int
	amount   int
	rand     *rand.Rand
}

func NewKeyParam(key Key) *KeyParam {
	source := rand.NewSource(time.Now().UnixNano())
	return &KeyParam{
		key:  key,
		rand: rand.New(source),
	}
}

func (b *KeyParam) Value() Key {
	if b.nextKey > 0 {
		return b.nextKey
	}
	return b.key
}

func (b *KeyParam) Last() Key {
	return b.lastKey
}

func (b *KeyParam) Display() string {
	return midi.Note(uint8(b.Value()))
}

func (b *KeyParam) Computed(root Key, scale Scale) Key {
	if b.nextKey > 0 {
		b.key = b.nextKey
		b.nextKey = 0
	}
	if b.amount == 0 {
		return b.key
	}
	key := Key(b.rand.Intn(int(math.Abs(float64(b.amount))) + 1))
	if b.amount > 0 {
		key = b.key + Key(key)
	} else {
		key = b.key - Key(key)
	}
	interval := key.AllSemitonesFrom(root)
	b.lastKey = b.key.Transpose(root, scale, interval)
	return b.lastKey
}

func (b *KeyParam) SetNext(key Key, root Key) {
	if key < minKey || key > maxKey {
		return
	}
	b.nextKey = key
	if int(key)+b.amount < int(maxKey) {
		b.amount++
	}
	if int(key)+b.amount > int(minKey) {
		b.amount--
	}
	b.interval = b.nextKey.AllSemitonesFrom(root)
}

func (b *KeyParam) Set(key Key) {
	b.key = key
	b.nextKey = 0
}

func (b *KeyParam) Transpose(root Key, scale Scale) {
	b.SetNext(b.key.Transpose(root, scale, b.interval), root)
}

func (b *KeyParam) RandomAmount() int {
	return b.amount
}

func (b *KeyParam) SetRandomAmount(amount int) {
	if int(b.key)+amount < int(minKey) || int(b.key)+amount > int(maxKey) {
		return
	}
	b.amount = amount
}

func (b *KeyParam) IsChanging() bool {
	return b.nextKey > 0
}

func (b *KeyParam) Name() string {
	return "key"
}

func (b *KeyParam) Symbol() string {
	return ""
}
