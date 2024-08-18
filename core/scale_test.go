package core

import "testing"

func TestScale(t *testing.T) {
	key := Key(63)
	root := Key(60)
	scale := IONIAN
	want := Key(64)

	if key.InScale(root, scale) {
		t.Fatalf("%d should not be in c ionian scale", key)
	}

	newKey := key.Transpose(root, scale, 3)
	if newKey != want {
		t.Fatalf("%d should transpose to %d in c ionian scale, got %d", key, want, newKey)
	}

	if !key.InScale(root, DORIAN) {
		t.Fatalf("%d should be in c dorian scale", key)
	}
}
