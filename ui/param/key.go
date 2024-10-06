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
}

func (k Key) Name() string {
	return k.nodes[0].(music.Audible).Note().Key.Name()
}

func (k Key) Display() string {
	if k.nodes[0].(music.Audible).Note().Key.RandomAmount() != 0 {
		return fmt.Sprintf(
			"%s%+d\u033c",
			k.nodes[0].(music.Audible).Note().Key.Display(),
			k.nodes[0].(music.Audible).Note().Key.RandomAmount(),
		)
	}
	return k.nodes[0].(music.Audible).Note().Key.Display()
}

func (k Key) Value() int {
	return int(k.nodes[0].(music.Audible).Note().Key.Value())
}

func (k Key) Increment() {
	k.Set(k.keyIndex() + 1)
}

func (k Key) Decrement() {
	k.Set(k.keyIndex() + -1)
}

func (k Key) Left() {
	k.SetAlt(k.nodes[0].(music.Audible).Note().Key.RandomAmount() - 1)
}

func (k Key) Right() {
	k.SetAlt(k.nodes[0].(music.Audible).Note().Key.RandomAmount() + 1)
}

func (k Key) Set(value int) {
	if value >= len(k.keys) {
		return
	}
	for _, n := range k.nodes {
		n.(music.Audible).Note().SetKey(k.keys[value], k.root)
	}
}

func (k Key) SetAlt(value int) {
	for _, n := range k.nodes {
		n.(music.Audible).Note().Key.SetRandomAmount(value)
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
		if k.nodes[0].(music.Audible).Note().Key.Value() == k.keys[i] {
			return i
		}
	}
	return 0
}
