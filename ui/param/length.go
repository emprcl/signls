package param

import (
	"fmt"
	"strconv"

	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
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
		return util.Normalize(
			fmt.Sprintf(
				"%s%+.1f\u033c",
				display,
				float64(l.nodes[0].(music.Audible).Note().Length.RandomAmount())/float64(pulsesPerStep),
			),
		)
	}
	return display
}

func (l Length) Value() int {
	return int(l.nodes[0].(music.Audible).Note().Length.Value())
}

func (l Length) AltValue() int {
	return 0
}

func (l Length) Up() {
	l.Set(l.Value() + 1)
}

func (l Length) Down() {
	l.Set(l.Value() - 1)
}

func (l Length) Left() {
	l.SetAlt(l.nodes[0].(music.Audible).Note().Length.RandomAmount() - 1)
}

func (l Length) Right() {
	l.SetAlt(l.nodes[0].(music.Audible).Note().Length.RandomAmount() + 1)
}

func (l Length) AltUp() {}

func (l Length) AltDown() {}

func (l Length) AltLeft() {}

func (l Length) AltRight() {}

func (l Length) Set(value int) {
	for _, n := range l.nodes {
		n.(music.Audible).Note().SetLength(uint8(value))
	}
}

func (l Length) SetAlt(value int) {
	for _, n := range l.nodes {
		n.(music.Audible).Note().Length.SetRandomAmount(value)
	}
}

func (l Length) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	l.Set(value)
}
