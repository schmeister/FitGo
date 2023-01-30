package funcs

import "math"

// Sine functions
func SineFunc(x float64, ps []float64) float64 {
	// amp, frequency, angle, decay
	return ps[0] * math.Sin(ps[1]*x+ps[2]) * math.Exp(-x*x*ps[3])
}

func SineDerv(x float64, ps []float64) (float64, float64) {
	// Find the wave that we are in.
	k := math.Trunc(x / math.Pi)

	// Only care about the frequency, not the amplitude or dampening.
	// Take the derivative of the sin function
	// Solve for zero (where the slope is zero)
	fPx := ps[0] * math.Cos(ps[1]*k+ps[2])
	resX := ((math.Pi / 2.0) + 2.0*k*math.Pi - ps[2]) / ps[1]
	return resX, fPx
}

// Polynomial functions
func PolyFunc(x float64, ps []float64) float64 {
	return ps[0]*x*x + ps[1]*x + ps[2]
}

func PolyDerv(_ float64, ps []float64) (float64, float64) {
	x := -ps[1] / (2 * ps[0])
	return x, PolyFunc(x, ps)
}