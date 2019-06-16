package plotter

type Range struct {
	Start int
	End   int
}

type Plotter struct {
	Xs Range
	Ys Range
	Zs Range
}

func FromString(plotStr string) Plotter {
	return Plotter{}
}

func (p Plotter) Eval(x, y, z int) (int, int, int) {
	return 0, 0, 0
}
