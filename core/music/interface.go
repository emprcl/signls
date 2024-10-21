package music

import "signls/core/common"

// Audible represents an interface for nodes that trigger notes.
type Audible interface {
	Arm()
	Note() *Note
	Muted() bool
	SetMute(mute bool)
	Trig(key Key, scale Scale, inDir common.Direction, pulse uint64)
	Emit(pulse uint64) []common.Direction
}
