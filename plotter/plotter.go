package plotter

import (
	"fmt"
)

// Parse takes the code definition of a plot, and outputs the list of coords
// which will be produced by running said plot
func Parse(plotStr string) ([]coord, error) {
	ast := parse(plotStr)
	if ast.Expressions.X == nil {
		return nil, fmt.Errorf("no expression for x coordinate given")
	}
	if ast.Expressions.Y == nil {
		return nil, fmt.Errorf("no expression for y coordinate given")
	}
	if ast.Expressions.Z == nil {
		return nil, fmt.Errorf("no expression for z coordinate given")
	}
	if len(ast.Ranges) == 0 {
		return nil, fmt.Errorf("no ranges given")
	}
	for i, r := range ast.Ranges {
		if r.Ident == nil || r.Start == nil || r.End == nil {
			return nil, fmt.Errorf("range %d is malformed", i+1)
		}
	}
	coords := ast.eval(prelude)
	return coords, nil
}
