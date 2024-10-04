package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
)

const (
	maxLength = 127
)

type Length struct {
	nodes []common.Node
}

func (l Length) Name() string {
	return "len"
}

func (l Length) Display() string {
	length := int(l.nodes[0].(music.Audible).Note().Length.Value())
	pulsesPerStep, stepsPerQuarterNote := l.nodes[0].(music.Audible).Note().ClockDivision()
	var display string
	switch length {
	case pulsesPerStep / 4:
		display = "1|64"
	case pulsesPerStep / 2:
		display = "1|32"
	case pulsesPerStep:
		display = "1|16"
	case pulsesPerStep * stepsPerQuarterNote / 2:
		display = "1|8"
	case pulsesPerStep * stepsPerQuarterNote:
		display = "1|4"
	case pulsesPerStep * stepsPerQuarterNote * 2:
		display = "1|2"
	case maxLength:
		display = "inf"
	default:
		display = fmt.Sprintf("%.1f", float64(length)/float64(pulsesPerStep))
	}
	if l.nodes[0].(music.Audible).Note().Length.RandomAmount() != 0 {
		return fmt.Sprintf(
			"%s%+.1f\u033c",
			display,
			float64(l.nodes[0].(music.Audible).Note().Length.RandomAmount())/float64(pulsesPerStep),
		)
	}
	return display
}

func (l Length) Value() int {
	return int(l.nodes[0].(music.Audible).Note().Length.Value())
}

func (l Length) Increment() {
	l.Set(l.Value() + 1)
}

func (l Length) Decrement() {
	l.Set(l.Value() - 1)
}

func (l Length) Left() {
	rand := l.nodes[0].(music.Audible).Note().Length
	rand.SetRandomAmount(rand.RandomAmount() - 1)
}

func (l Length) Right() {
	rand := l.nodes[0].(music.Audible).Note().Length
	rand.SetRandomAmount(rand.RandomAmount() + 1)
}

func (l Length) Set(value int) {
	for _, n := range l.nodes {
		n.(music.Audible).Note().SetLength(uint8(value))
	}
}
