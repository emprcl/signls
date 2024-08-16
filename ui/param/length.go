package param

import (
	"cykl/core"
	"fmt"
)

const (
	maxLength = 127
)

type Length struct {
	node core.Node
}

func (l Length) Name() string {
	return "len"
}

func (l Length) Display() string {
	length := int(l.node.(core.Emitter).Note().Length)
	pulsesPerStep, stepsPerQuarterNote := l.node.(core.Emitter).Note().ClockDivision()
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
	return int(l.node.(core.Emitter).Note().Length)
}

func (l Length) Increment() {
	l.Set(l.Value() + 1)
}

func (l Length) Decrement() {
	l.Set(l.Value() - 1)
}

func (l Length) Set(value int) {
	l.node.(core.Emitter).Note().SetLength(uint8(value))
}
