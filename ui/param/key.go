package param

import (
	"fmt"
	"time"

	"cykl/core/common"
	"cykl/core/music"
)

type Key struct {
	nodes []common.Node
	keys  []music.Key
	root  music.Key
	mode  KeyMode
}

func (k Key) Name() string {
	return k.mode.Display()
}

func (k Key) Display() string {
	switch k.nodes[0].(music.Audible).Note().Behavior.(type) {
	case music.SilentNote:
		return "•"
	case *music.RandomNote:
		return fmt.Sprintf(
			"%d oct. from %s",
			k.nodes[0].(music.Audible).Note().Behavior.Value(),
			k.nodes[0].(music.Audible).Note().KeyName(),
		)
	default:
		return k.nodes[0].(music.Audible).Note().KeyName()
	}
}

func (k Key) Value() int {
	return int(k.nodes[0].(music.Audible).Note().KeyValue())
}

func (k Key) Increment() {
	switch k.nodes[0].(music.Audible).Note().Behavior.(type) {
	case music.FixedNote:
		k.Set(k.keyIndex() + 1)
	default:
		k.mode.Increment()
	}
}

func (k Key) Decrement() {
	switch k.nodes[0].(music.Audible).Note().Behavior.(type) {
	case music.FixedNote:
		k.Set(k.keyIndex() + -1)
	default:
		k.mode.Decrement()
	}
}

func (k Key) Left() {
	k.mode.Left()
}

func (k Key) Right() {
	k.mode.Right()
}

func (k Key) Set(value int) {
	for _, n := range k.nodes {
		n.(music.Audible).Note().SetKey(k.keys[value], k.root)
	}
}

func (k Key) Preview() {
	go func() {
		n := *k.nodes[0].(music.Audible).Note()
		n.Play(music.Key(60), music.CHROMATIC)
		time.Sleep(300 * time.Millisecond)
		n.Stop()
	}()
}

func (k Key) keyIndex() int {
	for i := 0; i < len(k.keys); i++ {
		if k.nodes[0].(music.Audible).Note().KeyValue() == k.keys[i] {
			return i
		}
	}
	return 0
}
