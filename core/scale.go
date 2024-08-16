package core

type Scale uint16
type Interval uint16

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
			MINOR_7TH | MAJOR_7TH,
	)
	LYDIAN     = Scale(UNISON)
	MIXOLYDIAN = Scale(UNISON)
	AEOLIAN    = Scale(UNISON)
	LOCIRAN    = Scale(UNISON)
)

func (s Scale) Intervals() []int {
	intervals := []int{}
	for i := 0; i < 12; i++ {
		if s&(1<<i) != 0 {
			intervals = append(intervals, i)
		}
	}
	return intervals
}
