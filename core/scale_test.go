package core

import "testing"

func TestKeyTransposition(t *testing.T) {
	tests := []struct {
		key          Key
		previousRoot Key
		root         Key
		scale        Scale
		want         Key
	}{
		{Key(63), Key(60), Key(60), IONIAN, Key(64)},
		{Key(63), Key(60), Key(60), DORIAN, Key(63)},
		{Key(51), Key(60), Key(60), IONIAN, Key(52)},
		{Key(51), Key(60), Key(60), DORIAN, Key(51)},
		{Key(87), Key(60), Key(60), IONIAN, Key(88)},
		{Key(87), Key(60), Key(60), DORIAN, Key(87)},
		{Key(63), Key(60), Key(61), MIXOLYDIAN, Key(65)},
		{Key(68), Key(60), Key(61), MIXOLYDIAN, Key(70)},
	}
	for _, tt := range tests {
		oldInterval := tt.key.AllSemitonesFrom(tt.previousRoot)
		newKey := tt.key.Transpose(tt.root, tt.scale, oldInterval)
		if newKey != tt.want {
			t.Fatalf("%d should transpose to %d in %s %s scale, got %d", tt.key, tt.want, tt.root.Name(), tt.scale.Name(), newKey)
		}
	}
}
