package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/node"
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
	length := int(l.nodes[0].(*node.BaseEmitter).Note().Length)
	pulsesPerStep, stepsPerQuarterNote := l.nodes[0].(*node.BaseEmitter).Note().ClockDivision()
	switch length {
	case pulsesPerStep / 4:
		return "1|64"
	case pulsesPerStep / 2:
		return "1|32"
	case pulsesPerStep:
		return "1|16"
	case pulsesPerStep * stepsPerQuarterNote / 2:
		return "1|8"
	case pulsesPerStep * stepsPerQuarterNote:
		return "1|4"
	case pulsesPerStep * stepsPerQuarterNote * 2:
		return "1|2"
	case maxLength:
		return "inf"
	default:
		return fmt.Sprintf("%.1f", float64(length)/float64(pulsesPerStep))
	}
}

func (l Length) Value() int {
	return int(l.nodes[0].(*node.BaseEmitter).Note().Length)
}

func (l Length) Increment() {
	l.Set(l.Value() + 1)
}

func (l Length) Decrement() {
	l.Set(l.Value() - 1)
}

func (l Length) Left() {}

func (l Length) Right() {}

func (l Length) Set(value int) {
	for _, n := range l.nodes {
		n.(*node.BaseEmitter).Note().SetLength(uint8(value))
	}
}
