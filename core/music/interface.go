package music

import "cykl/core/common"

type Audible interface {
	Muted() bool
	SetMute(mute bool)
	Trig(key Key, scale Scale, inDir common.Direction, pulse uint64)
	Emit(pulse uint64) []common.Direction
}
