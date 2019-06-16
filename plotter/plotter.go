package plotter

// Parse takes the code definition of a plot, and outputs the list of coords
// which will be produced by running said plot
func Parse(plotStr string) []coord {
	ast := parse(plotStr)
	coords := ast.eval(nil)
	return coords
}
