package core

import (
	"maps"
	"math/bits"
	"slices"
)

const (
	UNISON Interval = 1 << iota
	MINOR_2ND
	MAJOR_2ND
	MINOR_3RD
	MAJOR_3RD
	FOURTH
	TRITONE
	FIFTH
	MINOR_6TH
	MAJOR_6TH
	MINOR_7TH
	MAJOR_7TH
)

const (
	CHROMATIC = Scale(
		UNISON | MINOR_2ND | MAJOR_2ND | MINOR_3RD | MAJOR_3RD |
			FOURTH | TRITONE | FIFTH | MINOR_6TH | MAJOR_6TH |
			MINOR_7TH | MAJOR_7TH,
	)
	IONIAN = Scale(
		UNISON | MAJOR_2ND | MAJOR_3RD |
			FOURTH | FIFTH | MAJOR_6TH |
			MAJOR_7TH,
	)
	DORIAN = Scale(
		UNISON | MAJOR_2ND | MINOR_3RD |
			FOURTH | FIFTH | MAJOR_6TH |
			MINOR_7TH,
	)
	PHRYGIAN = Scale(
		UNISON | MINOR_2ND | MINOR_3RD |
			FOURTH | FIFTH | MINOR_6TH |
			MINOR_7TH,
	)
	LYDIAN = Scale(
		UNISON | MAJOR_2ND | MAJOR_3RD |
			TRITONE | FIFTH | MAJOR_6TH |
			MAJOR_7TH,
	)
	MIXOLYDIAN = Scale(
		UNISON | MAJOR_2ND | MAJOR_3RD |
			FOURTH | FIFTH | MAJOR_6TH |
			MINOR_7TH,
	)
	AEOLIAN = Scale(
		UNISON | MAJOR_2ND | MINOR_3RD |
			FOURTH | FIFTH | MINOR_6TH |
			MINOR_7TH,
	)
	LOCRIAN = Scale(
		UNISON | MINOR_2ND | MINOR_3RD |
			FOURTH | TRITONE | MINOR_6TH |
			MINOR_7TH,
	)

	// TODO: ass more scale, pentatonic etc..
)

var (
	allScales = map[Scale]string{
		CHROMATIC:  "chromatic",
		IONIAN:     "ionian",
		DORIAN:     "dorian",
		PHRYGIAN:   "phrygian",
		LYDIAN:     "lydian",
		MIXOLYDIAN: "mixolydian",
		AEOLIAN:    "aeolian",
		LOCRIAN:    "locrian",
	}
)

type Interval uint16

func (i Interval) Semitones() int {
	return bits.TrailingZeros16(uint16(i))
}

type Scale uint16

func AllScales() []Scale {
	return slices.Collect(maps.Keys(allScales))
}

func (s Scale) Name() string {
	if name, ok := allScales[s]; ok {
		return name
	}
	return ""
}

func (s Scale) Intervals() []int {
	intervals := []int{}
	for i := 0; i < 12; i++ {
		if s&(1<<i) != 0 {
			intervals = append(intervals, i)
		}
	}
	return intervals
}
