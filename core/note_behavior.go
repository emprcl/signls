package core

import "math/rand"

type NoteBehavior interface {
	Name() string
	Play(n *Note, root Key, scale Scale)
}

func AllNoteBehaviors() []NoteBehavior {
	return []NoteBehavior{
		FixedNote{},
		SilentNote{},
		RandomNote{},
	}
}

type SilentNote struct{}

func (b SilentNote) Play(n *Note, root Key, scale Scale) {}
func (b SilentNote) Name() string {
	return "silent"
}

type FixedNote struct{}

func (b FixedNote) Play(n *Note, root Key, scale Scale) {
	if n.nextKey > 0 {
		n.Stop()
		n.Key = n.nextKey
	}
	n.Transpose(root, scale)
	n.midi.NoteOn(n.Channel, uint8(n.Key), n.Velocity)
}

func (b FixedNote) Name() string {
	return "fixed"
}

type RandomNote struct{}

func (b RandomNote) Play(n *Note, root Key, scale Scale) {
	n.Key = root + Key(rand.Intn(12))
	interval := n.Key.AllSemitonesFrom(root)
	n.Key = n.Key.Transpose(root, scale, interval)
	n.midi.NoteOn(n.Channel, uint8(n.Key), n.Velocity)
}

func (b RandomNote) Name() string {
	return "random"
}
