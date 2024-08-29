package core

import (
	"maps"
	"math"
	"slices"

	"cykl/midi"
)

// Interval constants represent various musical intervals as bitwise values.
// These intervals can be combined to form scales.
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

// Scale constants represent different musical scales as combinations of intervals.
const (
	CHROMATIC = Scale(UNISON | MINOR_2ND | MAJOR_2ND | MINOR_3RD | MAJOR_3RD | FOURTH |
		TRITONE | FIFTH | MINOR_6TH | MAJOR_6TH | MINOR_7TH | MAJOR_7TH)
	IONIAN     = Scale(UNISON | MAJOR_2ND | MAJOR_3RD | FOURTH | FIFTH | MAJOR_6TH | MAJOR_7TH)
	DORIAN     = Scale(UNISON | MAJOR_2ND | MINOR_3RD | FOURTH | FIFTH | MAJOR_6TH | MINOR_7TH)
	PHRYGIAN   = Scale(UNISON | MINOR_2ND | MINOR_3RD | FOURTH | FIFTH | MINOR_6TH | MINOR_7TH)
	LYDIAN     = Scale(UNISON | MAJOR_2ND | MAJOR_3RD | TRITONE | FIFTH | MAJOR_6TH | MAJOR_7TH)
	MIXOLYDIAN = Scale(UNISON | MAJOR_2ND | MAJOR_3RD | FOURTH | FIFTH | MAJOR_6TH | MINOR_7TH)
	AEOLIAN    = Scale(UNISON | MAJOR_2ND | MINOR_3RD | FOURTH | FIFTH | MINOR_6TH | MINOR_7TH)
	LOCRIAN    = Scale(UNISON | MINOR_2ND | MINOR_3RD | FOURTH | TRITONE | MINOR_6TH | MINOR_7TH)

	// TODO: add more scales, pentatonic, etc.
)

// allScales maps each scale constant to its corresponding name.
// This allows easy lookup of scale names based on their bitwise representation.
var allScales = map[Scale]string{
	CHROMATIC:  "chromatic",
	IONIAN:     "ionian",
	DORIAN:     "dorian",
	PHRYGIAN:   "phrygian",
	LYDIAN:     "lydian",
	MIXOLYDIAN: "mixolydian",
	AEOLIAN:    "aeolian",
	LOCRIAN:    "locrian",
}

// Key represents a musical note by its MIDI key number.
type Key uint8

// Name returns the name of the note according to the MIDI specification.
func (k Key) Name() string {
	return midi.Note(uint8(k))
}

// AllSemitonesFrom calculates the number of semitones between two keys.
func (k Key) AllSemitonesFrom(key Key) int {
	return int(k) - int(key)
}

// SemitonesFrom calculates the number of semitones between two keys,
// wrapped within an octave (0-11 semitones).
func (k Key) SemitonesFrom(key Key) int {
	d := int(k) - int(key)
	return mod(d, 12)
}

// InScale checks if the key is part of the given scale, relative to the root key.
func (k Key) InScale(root Key, scale Scale) bool {
	interval := k.SemitonesFrom(root)
	return scale&(1<<interval) != 0
}

// Transpose transposes the key according to the scale and old interval, adjusting
// the key if the new scale does not contain the exact transposition.
func (k Key) Transpose(root Key, scale Scale, oldInterval int) Key {
	// 1) First, let's just transpose a simple key change.
	newKey := Key(int(k) + oldInterval - k.AllSemitonesFrom(root))
	if newKey.InScale(root, scale) {
		return newKey
	}

	// 2) If not in scale, it means that the scale changed.
	// Let's try to change to push the key up or down
	// according to its initial interval, and check
	// if we're in the new scale.
	switch Interval(1 << (mod(oldInterval, 12))) {
	case MINOR_2ND, MINOR_3RD, MINOR_6TH, MINOR_7TH:
		newKey += Key(1)
	case MAJOR_2ND, MAJOR_3RD, MAJOR_6TH, MAJOR_7TH:
		newKey -= Key(1)
	case FOURTH:
		newKey += Key(1)
	case TRITONE, FIFTH:
		newKey -= Key(1)
	}

	if newKey.InScale(root, scale) {
		return newKey
	}

	// 3) If we're here, we're probably changing
	// to a scale with a different length.
	// ex: going from diatonic to pentatonic scale
	// Let's do best effort according to the min
	// distance to a note in the scale.
	// TODO: add test cases for this part.
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

// Interval represents a musical interval using a bitwise integer.
type Interval uint16

// Int returns the first interval set in the bitwise representation.
func (in Interval) Int() int {
	for i := 0; i < 12; i++ {
		if in&(1<<i) != 0 {
			return i
		}
	}
	return 0
}

// Scale represents a musical scale using a bitwise integer,
// where each bit corresponds to a semitone in an octave.
type Scale uint16

// AllScales returns a slice of all scales defined in the allScales map.
func AllScales() []Scale {
	return slices.Collect(maps.Keys(allScales))
}

// AllKeysInScale returns all MIDI keys within the given scale, relative to the root key.
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

// Name returns the name of the scale based on its bitwise representation.
func (s Scale) Name() string {
	if name, ok := allScales[s]; ok {
		return name
	}
	return ""
}

// Intervals returns a slice of intervals that make up the scale.
func (s Scale) Intervals() []int {
	intervals := []int{}
	for i := 0; i < 12; i++ {
		if s&(1<<i) != 0 {
			intervals = append(intervals, i)
		}
	}
	return intervals
}

// mod handles the modulo operation for negative numbers, ensuring
// the result is always non-negative.
func mod(a, b int) int {
	return (a%b + b) % b
}
