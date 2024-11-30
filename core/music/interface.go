package music

import (
	"signls/core/common"
	"signls/core/theory"
)

// Audible represents an interface for nodes that trigger notes.
type Audible interface {
	Arm()
	Note() *Note
	Muted() bool
	SetMute(mute bool)
	Trig(key theory.Key, scale theory.Scale, inDir common.Direction, pulse uint64)
	Emit(pulse uint64) []common.Direction
}
