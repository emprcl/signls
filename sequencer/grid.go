package sequencer

type Node interface {
	Emit(g *Grid)
	Trigger(g *Grid)
}

type Grid struct {
	nodes [][]Node
	h     int
	w     int
}

func (g *Grid) Compute() {
	g.RunEmitters()
	g.RunTriggers()
}

func (g *Grid) RunEmitters() {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			g.nodes[y][x].Emit(g)
		}
	}
}

func (g *Grid) RunTriggers() {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			g.nodes[y][x].Emit(g)
		}
	}
}
