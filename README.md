# FitGo: Finding the best RPM on a Swing Curve with GoLang, go-hep, and GoNum

I had been programming for quite some time, starting with an 8-bit Atari 400 with with only 8KB of RAM. My first professional task as a Software Engineer came many years later, but still over 30 years prior to 2023. I was in college working towards my BS in Computer science and had completed Physics I & II, Calculus I & II, and a number of other Computer Science pre-requisites. Modula-2, C, and the HP15c where my languages of choice. Ok, that last one is not a language, but a computation device - **the** Scientific Calculator of Engineers of the time. (One caveat, since this story happened over 30 years ago, many of the details have faded, I will do the best that I can to retell as accurate as I can.)

Recently being promoted to a Process Technician in the Photolithography department of a Semi-conductor Fabrication (FAB) plant. The opportunities to use my skills were going to be endless; the Photolighography area was ripe with areas for automation. One key part of the process is the coating of the silicon wafers with a photoresistive layer, dispensed from a package that could cover hundered, if not thousands of wafers. When the package needed to be replaced, many checks were in place to confirm that the process and quality was not altered. No matter how good the manufacturing of the resist was, there are still variations: The viscosity will change as well as the photoresitive components and many other attributes. The thickness of the resist needs to be extremely accurate and there are many engineering measurements and calculations necessary to ensure that. Obtaining the best spin speed (RPM) for the thickness is where my Software Engineering skill manifested itself which helped automate our processes and obtain the best (stable) configuration.

Resist Thickness is a key parameters that affects CD (Critical Dimension) in the lithography sequence. (https://ieeexplore.ieee.org/document/4529026)

The Swing Curve is a technique used that models the resist thickness, and incorporates thinfilm interferrence to determine the best thickness and stability. In general, this is a sine wave graph that is fit from a sampling of different resist thicknesses versus the CD size of the pattern being exposed on the wafers. This takes into account the very narrow light wavelength minimizing the affects from the numerous process variations.

![](https://github.com/schmeister/GoPlayground/blob/main/FitGo/testdata/Decay.png)
![](https://github.com/schmeister/FitGo/blob/main/testdata/SineWaveDecay.GIF)
<pre>
// Sine with decay
func SineFunc(x float64, ps []float64) float64 {
	// amp, frequency, angle, decay
	return ps[0] * math.Sin(ps[1]*x+ps[2]) * math.Exp(-x*x*ps[3])
}
</pre>

The dispensing and measurement phase is summarized as this: The resist is dispensed onto a series of wafers, with a DNS/Dainippon machine, each coated with a changing RPM, which resulted in a different thickness of resist. The thickness of each wafer is measured with a highly accurate measuring tools (https://www.kla.com/), exposed (https://www.asml.com/en), and developed, washing away the resist that was exposed. Finally the CDs were verified for optimal shape with a SEM (Scanning Electron Microscope). The raw Resist Thickness vs CD plotting would then look something like this:

![](https://github.com/schmeister/FitGo/blob/main/testdata/raw_poly.png)

Spin speed, temperatures, exposure rates, and a host of other variables may have micro-changes, the most stable location within the Swing Curve will needed to be found. This will be taken from the first derivative of the decaying sine wave formula. Remember, essentially the first derivative of a function results in the slope of that curve at a specific location.  In the case of our Sine curve, we are only concerned with part that defines the slope, and thus need to only take a partial derivative (the Sine part) and ignore the rest (the dampening part), though we will add the Amplitude back in just for appearances.

<pre>
func SineDerv(x float64, ps []float64) (float64, float64) {
	// Find the wave that we are in (approx)
	k := math.Trunc(x / math.Pi)

	// Only care about the frequency, not the amplitude or dampening.
	// Take the derivative of the sin function
	// Solve for zero (where the slope is zero)
	fPx := ps[0] * math.Cos(ps[1]*k+ps[2])
  
	// resx is the value within the wave we chose as our optimal.
	resX := ((math.Pi / 2.0) + 2.0*k*math.Pi - ps[2]) / ps[1]
	return resX, fPx
}
</pre>

## Step 1: Create some sample data with variability

<pre>
// Raw data
type Points struct {
	xdata []float64
	ydata []float64
}

// Raw data generation parameters
type Raw struct {
	Params    []float64
	XVariance float64
	YVariance float64
	Origin    float64
	Min       float64
	Max       float64
	Step      float64
	Steps     int16
	Func      MyFunc
}
....
	// raw data configuration, and Generate
	raw := Raw{
		// amp, omega, shift, decay
		Params:    []float64{10.0, 2.0, 0, .025},
		XVariance: 0.05,
		YVariance: 0.3,
		Min:       0,
		Max:       2,
		Step:      0.1,
		Steps:     10,
		Func:      funcs.SineFunc,
	}
	points := Generate(raw)

// Generate creates a number of points to attempt a fit to.
// There is a randomization added to each point to create some variability.
func Generate(raw Raw) Points {
	xdata := make([]float64, 0)
	ydata := make([]float64, 0)
	min := 0.0
	max := raw.Max - raw.Min
	step := (max - min) / float64(raw.Steps)
	for i := min; i <= max; i += step {
		// Create the real calculated value from the function f(x).
		val := raw.Func(i, raw.Params)

		// Include some Randomization and store in slice.
		xdata = append(xdata, i+raw.Min+(rand.Float64()*2*raw.XVariance))
		ydata = append(ydata, val+(rand.Float64()*2*raw.YVariance))
	}

	// Tell user what parameters were used to generate this data.
	fmt.Printf("generated = %.2f\n", raw.Params)

	return Points{xdata, ydata}
}
</pre>

## Step 2: Fit our formula to that data

From the raw data, we could probably come up with an approximation of the optimal spin speed, but that is not sufficient, we need as close as possible! The next task is to use those points and model them against the reference formula to find **the** best speed for our desired thickness.

Back when I originally wrote this, I used C and had to hand code my own minimizing function. Today, many of those tools are easily available and the actual implementation takes a significantly less amount of time. For my modern solution, using GoLang and minimizing and plotting functions makes this trivial.

<pre>
type Fitting struct {
	// Fit start data
	startParams []float64
	points      Points
	Func        MyFunc
}

...
	// sine parameters: [amp, omega, shift, decay]
	ExecuteFitAndGraph(points, funcs.SineFunc, funcs.SineDerv, []float64{10.0, 2.0, 0, .025}, "sine")
...
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
		plot_X_Label: "RPM (arbitrary)",
		plot_Y_Label: "CD (arbitrary Critical Dimension)",
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
</pre>

![](https://github.com/schmeister/FitGo/blob/main/testdata/fit_sine.png)

Looks like we got a pretty good fit!

## Step 3: Find the best (Zero) slope

I do the best RPM and plotting in the same section of code:

<pre>
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
	fitFunc MyFunc
	fitDerv MyDerv
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
</pre>

And the final result:

![](https://github.com/schmeister/FitGo/blob/main/testdata/best_sine.png)

## Step 4: but does it need to be this complicated?

Once we know the approximate location within the Swing Curve, we can simplify the fit formula and its derivative to a simple Polynomial function:

<pre>
// Polynomial functions
func PolyFunc(x float64, ps []float64) float64 {
	return ps[0]*x*x + ps[1]*x + ps[2]
}

// Derivative and solve for 0 - Zero slope! Best RPM!
func PolyDerv(_ float64, ps []float64) (float64, float64) {
	x := -ps[1] / (2 * ps[0])
	return x, PolyFunc(x, ps)
}
</pre>

![](https://github.com/schmeister/FitGo/blob/main/testdata/best_poly.png)

And the fit is only nominally different.
