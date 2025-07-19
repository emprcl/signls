package music

import (
	"errors"
	"regexp"
	"strconv"
)

var noteMap = map[string]int{
	"C":  0,
	"Db": 1,
	"D":  2,
	"Eb": 3,
	"E":  4,
	"F":  5,
	"Gb": 6,
	"G":  7,
	"Ab": 8,
	"A":  9,
	"Bb": 10,
	"B":  11,

	"c":  0,
	"db": 1,
	"d":  2,
	"eb": 3,
	"e":  4,
	"f":  5,
	"gb": 6,
	"g":  7,
	"ab": 8,
	"a":  9,
	"bb": 10,
	"b":  11,
}

// ConvertNoteToMIDI converts a musical note (ex A6 or Db3) to a midi note number
func ConvertNoteToMIDI(note string) (int, error) {
	re := regexp.MustCompile(`^([A-Ga-g][b]?)(-?\d+)$`)
	matches := re.FindStringSubmatch(note)
	if len(matches) < 3 {
		return 0, errors.New("invalid note format")
	}

	noteName := matches[1]
	octaveStr := matches[2]

	octave, err := strconv.Atoi(octaveStr)
	if err != nil {
		return 0, errors.New("invalid octave")
	}

	if base, ok := noteMap[noteName]; ok {
		return base + octave*12, nil
	}
	return 0, errors.New("unknown note name")
}
