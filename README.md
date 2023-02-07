# FitGo: Finding the best RPM on a Swing Curve with GoLang, go-hep, and GoNum

In my early 20's I had my first task as a professional software engineer. I was in college working towards my BS in Computer Science and had completed Physics I & II, Calculus I & II, and a number of undergraduate Computer Science courses. Modula-2, C, and the HP-15C were my tools of choice. Ok, that last one is not count, but it was **the** Scientific Calculator of Engineers at the time. I had been programming for quite some time prior to this event. Years before, I started my interest in software by coding on an 8-bit Atari 400 with only 8KB of RAM. One of my first memories of programming a computer is typing the full issue of Antic V1-N5 (Dec. 1982 - https://archive.org/details/1982-12-anticmagazine/mode/2up) into my Atari, I would do this over and over (we had no non-volatile memory - disk, tape, etc.) modifying the code each time to see what happened (my favorite was Bats - page 57-60). I loved to experiment with the code modifying the physics and controls until I had a game that was uniquely mine. I took pages and pages of notes and probably could have written it from scratch at that point.

**Warning**: the following events took place over 30 years ago, many of the details have faded, as I am also sure, many of the technologies and processes have changed.

Fast forward a few years, I had recently been promoted to a Process Technician in the photolithography department of a semi-conductor fabrication facility (FAB). The opportunities to use my skills were going to be endless. The Photolighography area was ripe for computer automation and only computers could efficiently analyze the magnitude of data being generated.

One key part of the process is the coating of the silicon wafers with a photoresistive layer. The dispensing mechanism was extremely accurate, but no matter how accurate the tools were, slight variations do occur. These results of these slight variations need to be minimized when working with incredibly small features as where required for semi-conductors. As an example, when the resist package was replaced, many checks were required to confirm that the resist thinfilm coating was still within the specifications. As mentioned before, no matter how consistent the manufacturing of the resist was, there are still variations: viscosity, refractive index, photosensitive dyes, and many other attributes. The thickness of the resist needs to be extremely accurate and is one of the most important aspects of semiconductor manufacturing. (https://ieeexplore.ieee.org/document/4529026). 

Determining the optimal spin speed (RPM) to obtain the best thickness is where my software engineering skills manifested. The process at the time would require the engineer doing the checks to break out their HP 15C, enter the data and determine the best RPM. A fairly slow and error prone process.

The Swing Curve was a technique used that models the resist thickness, taking into account thinfilm interferrence to determine the most stable point. In brief, this formula is a decaying sine wave graph that is fit from a sampling of different thicknesses versus the CD size (Critical Dimension - feature size) of the pattern being exposed on the wafers. As the wafers are exposed with a very narrow spectrum of light, constructive and destructive interferrence **could be a big problem**.

![](https://github.com/schmeister/FitGo/blob/main/testdata/Decay.png)
![](https://github.com/schmeister/FitGo/blob/main/testdata/SineWaveDecay.GIF)
<pre>
// Sine with decay
func SineFunc(x float64, ps []float64) float64 {
	// amp, frequency, angle, decay
	return ps[0] * math.Sin(ps[1]*x+ps[2]) * math.Exp(-x*x*ps[3])
}
</pre>

The dispensing and measurement phase is summarized as this: The resist is dispensed onto a series of wafers, each coated with a changing RPM, which resulted in a different thickness of resist. The thickness of the resist is measured for each RPM with a highly accurate measuring tools (https://www.kla.com/), exposed (https://www.asml.com/en), and developed, which "washes" away the resist that was exposed (positive resist). Finally the CDs were verified for optimal shape with a SEM (Scanning Electron Microscope). The Resist Thickness vs CD was plotted and the raw data may look something like this:

![](https://github.com/schmeister/FitGo/blob/main/testdata/raw_poly.png)

As previously mentioned, spin speed, temperatures, exposure rates, and other variables may cause micro-changes in the thickness, the most stable location within the Swing Curve will needed to be found. This "optimal RPM" will be located by using the first derivative of the decaying sine wave formula. Remember, essentially the first derivative of a function results in the slope of the original formula at the given point, a slope of Zero (0) is optimal in this situation. Taking the derivative, in the case of our sine curve, we are only concerned with the part that defines the slope, and thus only need to take a partial derivative of the frequency (the sine part), and can ignore the the dampening/decaying part.

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

This randomized data will be used as input to our system to check that the fit and best RPM is indeed found.  See graph above with plotted raw data.

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
....
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

From visual observation of the raw data, we could probably come up with an approximation of an acceptable RPM, but that is not sufficient. The slightest variation in any of the steps would potentially cause a significant amount of rework if the zero slope is not chosen. We can do better than the approximation. The task is to use those raw data points and model them against the reference formula to find **the** best speed to achieve our desired thickness. To do this, we need to fit the curve to those raw points. 

Back when I wrote my original version, the techstack consisted of a VAX/VMS system with RS/1 (https://www.jstor.org/stable/1309968), C, and a hand coded (by myself) Least Squares minimizing function. Nearly no pre-written frameworks to use, and those available required a significant amount of code to utilize properly. Today, many of those tools are easily available and actual implementation takes a significantly less amount of time. For my modern reincarnation, as you can tell, I am using GoLang and import minimizing and plotting packages making this whole solution quite trivial. As a rough order of magnitude, I would say it was 20 lines of code back then to every 1 line of code currently.

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

## Step 3: Find the optimal slope

In my current version, I do the best RPM calulations and plotting in the same section of code. The best RPM is provided in the Key, and the zero slope line is attached to the sine curve at the best location. 

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

Once we know the approximate location within the Swing Curve, we can simplify the fit formula and its derivative to a simple Polynomial function. The use of the Decaying Sine function need only be applied when experimenting with larger ranges of thicknesses. Thickness variations that may span multiple waves.

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

Wow, that formula and derivation is significantly more straight forward!

![](https://github.com/schmeister/FitGo/blob/main/testdata/best_poly.png)

And the best fit from the Sine vs Polynomial is only nominally different. Though using a computer, the hardest part may be determining the proper derivative for the complex forumula.
