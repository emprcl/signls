package param

import (
	"fmt"
	"time"

	"cykl/core/common"
	"cykl/core/music"
)

type KeyMode uint8

const (
	KeyModeRandom KeyMode = iota
	KeyModeSilent
)

type Key struct {
	nodes []common.Node
	keys  []music.Key
	root  music.Key
	mode  KeyMode
}

func (k *Key) Name() string {
	return k.nodes[0].(music.Audible).Note().Key.Name()
}

func (k *Key) Display() string {
	if k.mode == KeyModeSilent {
		return "⨯"
	}

	if k.nodes[0].(music.Audible).Note().Key.RandomAmount() != 0 {
		return fmt.Sprintf(
			"%s%+d\u033c",
			k.nodes[0].(music.Audible).Note().Key.Display(),
			k.nodes[0].(music.Audible).Note().Key.RandomAmount(),
		)
	}
	return k.nodes[0].(music.Audible).Note().Key.Display()
}

func (k *Key) Value() int {
	return int(k.nodes[0].(music.Audible).Note().Key.Value())
}

func (k *Key) AltValue() int {
	switch k.mode {
	case KeyModeSilent:
		return 0
	default:
		return k.nodes[0].(music.Audible).Note().Key.RandomAmount()
	}
}

func (k *Key) Up() {
	k.Set(k.keyIndex() + 1)
}

func (k *Key) Down() {
	k.Set(k.keyIndex() - 1)
}

func (k *Key) Left() {
	k.SetAlt(k.AltValue() - 1)
}

func (k *Key) Right() {
	k.SetAlt(k.AltValue() + 1)
}

func (k *Key) AltUp() {}

func (k *Key) AltDown() {}

func (k *Key) AltLeft() {
	k.mode = KeyMode((k.mode - 1) % 2)
	for _, n := range k.nodes {
		n.(music.Audible).Note().Key.SetSilent(k.mode == KeyModeSilent)
	}
}

func (k *Key) AltRight() {
	k.mode = KeyMode((k.mode + 1) % 2)
	for _, n := range k.nodes {
		n.(music.Audible).Note().Key.SetSilent(k.mode == KeyModeSilent)
	}
}

func (k *Key) Set(value int) {
	if k.mode == KeyModeSilent {
		return
	}
	if value < 0 || value >= len(k.keys) {
		return
	}
	for _, n := range k.nodes {
		n.(music.Audible).Note().SetKey(k.keys[value], k.root)
	}
}

func (k *Key) SetAlt(value int) {
	switch k.mode {
	case KeyModeSilent:
		return
	default:
		for _, n := range k.nodes {
			n.(music.Audible).Note().Key.SetRandomAmount(value)
		}
	}
}

func (k *Key) Preview() {
	go func() {
		n := *k.nodes[0].(music.Audible).Note()
		n.Play()
		time.Sleep(300 * time.Millisecond)
		n.Stop()
	}()
}

func (k *Key) keyIndex() int {
	for i := 0; i < len(k.keys); i++ {
		if k.nodes[0].(music.Audible).Note().Key.Value() == k.keys[i] {
			return i
		}
	}
	return 0
}
