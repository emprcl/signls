package music

import "cykl/midi"

const (
	maxKey Key = 127
)

type KeyParam struct {
	key      Key
	nextKey  Key
	interval int
}

func (b *KeyParam) Value() Key {
	if b.nextKey > 0 {
		return b.nextKey
	}
	return b.key
}

func (b *KeyParam) Display() string {
	return midi.Note(uint8(b.Value()))
}

func (b *KeyParam) Computed() Key {
	if b.nextKey > 0 {
		b.key = b.nextKey
		b.nextKey = 0
	}
	return b.key
}

func (b *KeyParam) SetNext(key Key, root Key) {
	if key > maxKey {
		b.nextKey = Key(0)
	} else {
		b.nextKey = key
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

func (b *KeyParam) IsChanging() bool {
	return b.nextKey > 0
}

func (b *KeyParam) Name() string {
	return "key"
}

func (b *KeyParam) Symbol() string {
	return ""
}
