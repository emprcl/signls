package core

import (
	"cykl/midi"
	"fmt"
	"testing"
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
	for _, v := range benchmarks {
		b.Run(fmt.Sprintf("grid_size_%dx%d", v.size, v.size), func(b *testing.B) {
			grid := NewGrid(v.size, v.size, &midi.Mock{})
			for i := 0; i < b.N; i++ {
				grid.Update()
			}
		})
	}
}
