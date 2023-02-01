package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"github.com/schmeister/FitGo/pkg/funcs"
)

type FitFunc func(x float64, ps []float64) float64
type FitDerv func(x float64, ps []float64) (float64, float64)

type Points struct {
	xdata []float64
	ydata []float64
}

// Raw data generation
type Raw struct {
	Params    []float64 // Parameters used in the FitFunc
	Func      FitFunc	// Function used to generate the raw data
	XVariance float64	// Amount of variance in X direction
	YVariance float64	// Amount of variance in Y direction
	Origin    float64
	XMin       float64
	YMax       float64
	Step      float64
	Steps     int16
}

type Fitting struct {
	// Fit start data
	startParams []float64
	points      Points
	Func        FitFunc
}

type Plotting struct {
	fitted []float64
	raw    Points

	// Graph display
	plot_X_Label string
	plot_Y_Label string
	plot_X_Min   float64
	plot_X_Max   float64
	plot_Y_Min   float64
	plot_Y_Max   float64

	// Formula
	fitFunc FitFunc
	fitDerv FitDerv
}

// Main
func main() {
	rand.Seed(time.Now().UnixNano())

	// Define raw data generator.
	raw := Raw{
		// amp, omega, shift, decay
		Params:    []float64{10.0, 2.0, 0, .025},
		Func:      funcs.SineFunc,
		XVariance: 0.05,
		YVariance: 0.3,
		XMin:       0,
		YMax:       2,
		Step:      0.1,
		Steps:     10,
	}
	points := raw.Generate()

	// sine parameters: [amp, omega, shift, decay]
	ExecuteFitAndGraph(points, funcs.SineFunc, funcs.SineDerv, []float64{10.0, 2.0, 0, .025}, "sine")
	ExecuteFitAndGraph(points, funcs.PolyFunc, funcs.PolyDerv, []float64{-10.0, 15.0, 5}, "poly")
}

// Generate creates a number of points to attempt a fit to.
// There is a randomization added to each point to create some variability.
func (raw Raw) Generate() Points {
	xdata := make([]float64, 0)
	ydata := make([]float64, 0)
	min := 0.0
	max := raw.YMax - raw.XMin
	step := (max - min) / float64(raw.Steps)
	for i := min; i <= max; i += step {
		// Create the real calculated value from the function f(x).
		val := raw.Func(i, raw.Params)

		// Include some Randomization and store in slice.
		xdata = append(xdata, i+raw.XMin+(rand.Float64()*2*raw.XVariance))
		ydata = append(ydata, val+(rand.Float64()*2*raw.YVariance))
	}

	// Tell user what parameters were used to generate this data.
	fmt.Printf("generated = %.2f\n", raw.Params)

	return Points{xdata, ydata}
}

// ExecuteFitAndGraph takes in the raw points, and models the data to the fit function.
// Then graphs both the fitted function and the derivative function (get get the slope).
func ExecuteFitAndGraph(points Points,
	fit func(float64, []float64) float64,
	derv func(float64, []float64) (float64, float64),
	fitParams []float64,
	name string) {

	// Attempt to fit with these starting parameters.
	// These parameters may need to be changed based on the formula and data.
	fitting := Fitting{
		startParams: fitParams,
		points:      points,
		Func:        fit,
	}
	fitted := fitting.Fit()
	fmt.Printf("got (%4s)= %.2f\n", name, fitted)

	// Data is attempted to be fitted to, time to plot all data.
	plotting := Plotting{
		fitted:       fitted,
		raw:          points,
		plot_X_Label: "some arbitrary RPM",
		plot_Y_Label: "arbitrary CD (Critical Dimension)",
		plot_X_Min:   points.xdata[0],
		plot_X_Max:   points.xdata[len(points.xdata)-1],
		plot_Y_Min:   -10.0,
		plot_Y_Max:   13.0,
		fitFunc:      fit,
		fitDerv:      derv,
	}

	// Plot the same graph 3 times each showing more data.
	plotting.PlotData(true, false, false, "raw_"+name+".png")
	plotting.PlotData(true, true, false, "fit_"+name+".png")
	plotting.PlotData(true, true, true, "best_"+name+".png")
}

// Fit uses the fitting parameters to attempt the fit.
// The return number are those that should/could be used in the
// given function.
func (fitting *Fitting) Fit() []float64 {
	// Fit to Raw data
	res, err := fit.Curve1D(
		fit.Func1D{
			F:  fitting.Func,
			X:  fitting.points.xdata,
			Y:  fitting.points.ydata,
			Ps: fitting.startParams,
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	return res.X
}

// PlotData does just that plots all the data and then at the same time
// uses the Dervative function to plot and compute the best RPM.
func (plotting *Plotting) PlotData(showRaw, showFit, showSlope bool, name string) {
	p := hplot.New()
	p.X.Label.Text = plotting.plot_X_Label
	p.Y.Label.Text = plotting.plot_Y_Label
	p.X.Min = plotting.plot_X_Min
	p.X.Max = plotting.plot_X_Max
	p.Y.Min = plotting.plot_Y_Min
	p.Y.Max = plotting.plot_Y_Max

	// Raw Data - if requested
	if showRaw {
		s := hplot.NewS2D(hplot.ZipXY(plotting.raw.xdata, plotting.raw.ydata))
		s.Color = color.RGBA{0, 0, 255, 255}
		p.Add(s)
		p.Legend.Add("Raw", s)
	}

	// Fit Function - if requested
	if showFit {
		f := plotter.NewFunction(func(x float64) float64 {
			y := plotting.fitFunc(x, plotting.fitted)
			return y
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 100
		p.Add(f)
		p.Legend.Add("Fit", f)

	}

	// Derivative - if requested
	if showSlope {
		e := plotter.NewFunction(func(x float64) float64 {
			_, y := plotting.fitDerv(x, plotting.fitted)
			return y
		})
		e.Color = color.RGBA{0, 255, 0, 255}
		e.Samples = 100
		e.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		p.Add(e)

		// The follow median may not result in the best location for a Sine wave fit.
		x, _ := plotting.fitDerv((p.X.Min+p.X.Max)/2, plotting.fitted)
		best := fmt.Sprintf("Best: %.2f", x)
		p.Legend.Add(best, e)
	}

	// Background Grid
	p.Add(plotter.NewGrid())

	// Set Legend Position
	p.Legend.Left = true
	p.Legend.Top = true

	// Save file
	err := p.Save(20*vg.Centimeter, -1, "testdata/"+name)
	if err != nil {
		log.Fatal(err)
	}
}
