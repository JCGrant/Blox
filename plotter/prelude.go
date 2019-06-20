package plotter

import "math"

var prelude = map[string]interface{}{
	"sin": oneArg(math.Sin),
	"cos": oneArg(math.Cos),
}

func noArgs(fn func() float64) func(...float64) float64 {
	return func(fs ...float64) float64 {
		return fn()
	}
}

func oneArg(fn func(float64) float64) func(...float64) float64 {
	return func(fs ...float64) float64 {
		return fn(fs[0])
	}
}

func twoArgs(fn func(float64, float64) float64) func(...float64) float64 {
	return func(fs ...float64) float64 {
		return fn(fs[0], fs[1])
	}
}
