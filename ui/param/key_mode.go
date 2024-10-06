package param

// import (
// 	"cykl/core/common"
// 	"cykl/core/music"
// )

// type KeyMode struct {
// 	nodes []common.Node
// 	modes []music.NoteBehavior
// }

// func (k KeyMode) Name() string {
// 	return "mode"
// }

// func (k KeyMode) Display() string {
// 	return k.nodes[0].(music.Audible).Note().Behavior.Name()
// }

// func (k KeyMode) Value() int {
// 	return 0
// }

// func (k KeyMode) Increment() {
// 	val := k.nodes[0].(music.Audible).Note().Behavior.Value() + 1
// 	k.nodes[0].(music.Audible).Note().Behavior.Set(val)
// }

// func (k KeyMode) Decrement() {
// 	val := k.nodes[0].(music.Audible).Note().Behavior.Value() - 1
// 	k.nodes[0].(music.Audible).Note().Behavior.Set(val)
// }

// func (k KeyMode) Left() {
// 	k.Set(k.keyModeIndex() - 1)
// }

// func (k KeyMode) Right() {
// 	k.Set(k.keyModeIndex() + 1)
// }

// func (k KeyMode) Set(value int) {
// 	if value < 0 {
// 		value = len(k.modes) - 1
// 	} else if value >= len(k.modes) {
// 		value = 0
// 	}
// 	for _, n := range k.nodes {
// 		n.(music.Audible).Note().Behavior = k.modes[value]
// 	}
// }

// func (k KeyMode) keyModeIndex() int {
// 	for i := 0; i < len(k.modes); i++ {
// 		if k.nodes[0].(music.Audible).Note().Behavior == k.modes[i] {
// 			return i
// 		}
// 	}
// 	return 0
// }
