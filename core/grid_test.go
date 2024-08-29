package core

import (
	"fmt"
	"testing"

	"cykl/midi"
)

var benchmarks = []struct {
	size int
}{
	{size: 20},
	{size: 100},
	{size: 200},
	{size: 300},
	{size: 400},
	{size: 500},
}

func BenchmarkGrid(b *testing.B) {
	midi := &midi.Mock{}
	for _, v := range benchmarks {
		b.Run(fmt.Sprintf("grid_size_%dx%d", v.size, v.size), func(b *testing.B) {
			grid := NewGrid(v.size, v.size, midi)
			grid.AddNode(NewBangEmitter(midi, DOWN|RIGHT, true), 7, 7)
			grid.AddNode(NewSpreadEmitter(midi, DOWN), 11, 7)
			grid.AddNode(NewSpreadEmitter(midi, LEFT), 11, 11)
			grid.AddNode(NewSpreadEmitter(midi, UP), 7, 11)
			grid.AddNode(NewBangEmitter(midi, RIGHT, true), 7, 2)
			grid.AddNode(NewSpreadEmitter(midi, LEFT), 12, 2)
			grid.AddNode(NewBangEmitter(midi, RIGHT, true), 7, 3)
			grid.AddNode(NewSpreadEmitter(midi, LEFT), 9, 3)
			for i := 0; i < b.N; i++ {
				grid.Update()
			}
		})
	}
}
