package plotter

import "math"

var prelude = map[string]interface{}{
	"sin": func(fs ...float64) float64 {
		return math.Sin(fs[0])
	},
	"cos": func(fs ...float64) float64 {
		return math.Cos(fs[0])
	},
}
