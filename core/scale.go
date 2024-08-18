package core

import (
	"cykl/midi"
	"maps"
	"math"
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

type Key uint8

func (k Key) Name() string {
	return midi.Note(uint8(k))
}

func (k Key) AllSemitonesFrom(key Key) int {
	return int(key) - int(k)
}

func (k Key) SemitonesFrom(key Key) int {
	d := int(key) - int(k)
	if d < 0 {
		d = -d
	}
	return d % 12
}

func (k Key) InScale(root Key, scale Scale) bool {
	interval := k.SemitonesFrom(root)
	return scale&(1<<interval) != 0
}

func (k Key) Transpose(root Key, scale Scale, oldInterval int) Key {
	// 1) First, lets just transpose a simple key change.
	newKey := Key(int(k) + k.AllSemitonesFrom(root) - oldInterval)
	if newKey.InScale(root, scale) {
		return newKey
	}

	// 2) If not in scale, it means that the scale changed.
	// Lets try to change to push the key up or down
	// according to its initial interval, and check
	// if we're in the new scale.'
	if oldInterval < 0 {
		oldInterval = -oldInterval
	}
	switch Interval(1 << (oldInterval % 12)) {
	case MINOR_2ND, MINOR_3RD, MINOR_6TH, MINOR_7TH:
		newKey += Key(1)
	case MAJOR_2ND, MAJOR_3RD, MAJOR_6TH, MAJOR_7TH:
		newKey -= Key(1)
	case FOURTH:
		newKey += Key(1)
	case TRITONE:
		newKey -= Key(1)
	case FIFTH:
		newKey -= Key(1)
	}

	if newKey.InScale(root, scale) {
		return newKey
	}

	// 3) If we're here, we're probably changing
	// to a scale with a different length.
	// ex: going from diatonic to pentatonic scale
	// Let's do best effort according to the min
	// distance to a note in the scale.'
	minDistance := math.MaxUint8
	interval := newKey.SemitonesFrom(root)
	for i := 0; i < 12; i++ {
		if scale&(1<<i) != 0 {
			distance := interval - i
			if distance < 0 {
				distance = -distance
			}
			if distance < minDistance {
				newKey = root + Key(i)
				minDistance = distance
			}
		}
	}
	return newKey
}

type Interval uint16

func (in Interval) Int() int {
	for i := 0; i < 12; i++ {
		if in&(1<<i) != 0 {
			return i
		}
	}
	return 0
}

type Scale uint16

func AllScales() []Scale {
	return slices.Collect(maps.Keys(allScales))
}

func AllKeysInScale(root Key, scale Scale) []Key {
	var keys []Key
	for i := 0; i <= 127; i++ {
		if scale&(1<<(i%12)) != 0 {
			key := root%12 + Key(i)
			keys = append(keys, key)
		}
	}
	return keys
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
