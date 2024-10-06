package music

type KeyBehavior interface {
	Value() Key
	Computed() Key
	Set(key Key)
	SetNext(key Key, root Key)
	IsChanging() bool
	Transpose(root Key, scale Scale)
	Name() string
	Symbol() string
}

type FixedKey struct {
	key      Key
	nextKey  Key
	interval int
}

func (b *FixedKey) Value() Key {
	if b.nextKey > 0 {
		return b.nextKey
	}
	return b.key
}

func (b *FixedKey) Computed() Key {
	if b.nextKey > 0 {
		b.key = b.nextKey
	}
	return b.key
}

func (b *FixedKey) SetNext(key Key, root Key) {
	if key > maxKey {
		b.nextKey = Key(0)
	} else {
		b.nextKey = key
	}
	b.interval = b.nextKey.AllSemitonesFrom(root)
}

func (b *FixedKey) Set(key Key) {
	b.key = key
	b.nextKey = 0
}

func (b *FixedKey) Transpose(root Key, scale Scale) {
	b.key = b.key.Transpose(root, scale, b.interval)
}

func (b *FixedKey) IsChanging() bool {
	return b.nextKey > 0
}

func (b *FixedKey) Name() string {
	return "key"
}

func (b *FixedKey) Symbol() string {
	return ""
}
