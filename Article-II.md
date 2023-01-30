# Solving a Swing Curve with GO and go-hep and GoNum

My first professional task as a Software Engineer came just over 30 years ago. I was in college working towards my BS in Computer science and had completed Physics I & II, Calculus I & II, and a number of other Computer Science pre-requisites. Modula-2, C, and HP15c where my languages of choice. Ok, that last one is not a language, but a computation device - **the** Scientific Calculator of Engineers.

I had recently been promoted to a Process Technician in the Photolithography department of a Semi-conductor Fabrication (FAB) plant. The photolithography department is the part of the FAB that actually places the patterns on a thinfilm prior to etching. Part of this process is the replacement of a photosensitive layer, known as a resist. The thickness of the resist needs to be very accurate and there are many engineering processes necessary to ensure the proper thickness. Obtaining the best parameters for the thickness is where my Software Engineering skill manifested itself which helped automate our processes.

The Swing Curve is a technique that models the resist thickness, and incorporates thinfilm interferrence to determine the best thickness of the thinfile. In general, this is a sine wave graph that is fit from a sampling of different resist thicknesses versus the Critical Dimension (CD) size of the pattern being exposed on the wafers. 
**Show a swing curve with some offset and decay**

The resist is dispensed on a series of wafers, with a DNS/Dainippon tool, each with a different resist thickness, determined by the spin speed during dispensing. The thickness was measured with highly accurate measuring tools (https://www.kla.com/), exposed (https://www.asml.com/en), and finally the CDs were verified with a SEM (Scanning Electron Microscope). The raw CD vs Thickness step would look something like this:
**Insert raw data graph**

Spin speed, temperatures, exposure rates, and a host of other variables may have micro-changes, the most stable location within the Swing Curve needed to be found (https://ieeexplore.ieee.org/document/4529026 and https://en.wikipedia.org/wiki/Thin-film_interference). From the raw data, we could probably come up with an approximation of the optimal spin speed, but that is not sufficient, we need as close as possible! The next task is to use those points and model them against a known formula to find **the** best speed for our desired thickness.

A sine wave with a decay was one of the most common formulas that could be used: f(x) = Amplitude * sin(Freq * x + Shift) * e^((-x^2)*Decay).

Simple enough, using a minimizing function we can find the best fit by modifying the Amplitude, Lambda, Shift, and Decay rate. The next step would be to find were the RPM where the slope is **Zero**, the most error proof thickness is.

resX := ((math.Pi / 2.0) + 2.0*k*math.Pi - ps[2]) / ps[1]
now to find the optimal slope, let's compute the derivative, but we don't need the full derivative, we only need to determine when the sine function has zero slope, and this is the first derviative of f(x): f'(x) = cos(Lambda * x + Shift) = 0. By setting f'(x) to 0, that means we are going to solve f'(x) where f'(x) == 0.

The first task is to fit the formula to the raw data.
**Show and describe FIT code to do this"
**Insert graph with fit curve**

Great, we have modelled a function to the raw data. Next step take the derivative:
**Show and describe Derivative code to do this"
**Insert graph with derivative curve**

What value do we chose from the derivation curve? As we previously said, we want to choose where f'(x) == 0. So choose where the sinusoidal curve is right where y==0.

If you notice, in this example I used a very accurate Sine function to solve for the best spin speed. But we only used one peak and the slopes on either side. There is another much simpiler formula, though not suitable for wide ranges, that we can use since we should already have a baseline for our initial spin speed. You guessed it - a polynomial function:  f(x) = anxn + ... + a1x + a0

//////////////////////

I read alot, most technical articles, but occationally a random fiction book. One of the topics in the articles I see quite often are those saying that "software engineers" are the same as "software developers" and are the same as "programers". I am going to throw my opinions into this mix in a technical sort of way.

These events occured over 30 years ago, so many of the finer details have been lost to me over time, but I also have a few reference links at the end if you would like to do more research. My first real job was working as an operator within a clean room, specifically the photolithography department imaging circuits on silicon wafers. After a few years of that I was partially through my Bachelor'd degree in Computer Science. I had completed many of my Science Track and math courses. Calculus I & II and Physics I & II were under my belt. I was unbeatable. I was promoted to Process Technician, someone that helped do the common higher level engineering tasks to keep the department running efficiently.

One of the common tasks within the photolithography department was the changing of the Photoresist on the wafer coating machines. The photoresist is carefully dispensed and spun onto the wafers with very precise equipment. With the new resist, a number of checks need to be performed. The remainder of this article will discuss how the system is calibrated to make sure the resist thickness is precisely what is required.

Resist Thickness needs to adhere to a certain few standards and the main topic we will look at:
1) The Resist must be between a minimum and maximum thickness.
2) The thickness needs to be such that the process is in the most stable state possible, meaning that variations in exposure and thickness result is the least affects. (The photo-resist is exposed with a specific wave-length of light, which may cause constructive and destructive interferrence)

Skipping many other steps, thickness data was obtained that would be used to generate a "swing curve". With the **swing curve**, we could determine the best spin speed and exposure for this batch of resist. A key and very important step to making all silicon basid chips.

The initial data may look somethink like this:
 **Insert raw data graph**

Many of the engineer would break out their slide rules (well not really, just their HP15c calculators) to compute the best spin speed. What I did was automate this process by writing an application in C that automated all the calculations. 

The result would look similar to this, fitting a function to the raw data, then allowing us to get the exact best spin speed.
**Insert raw data and fitted function**





-------------


The steps taken with new resist was roughly as follows. A sample number of wafers are coated, each with a slightly different spin speed. Each wafer is then measured to get the precise thickness of the resist. Each wafer is exposed, developed (which removes the exposed resist) then the final step was to use a Scanning Electron Microscope to see finer details.

References:
https://snfexfab.stanford.edu/guide/equipment/svg-resist-coat-track-2-svgcoat2
https://www.kla.com/products/instruments/thin-film-reflectometers
https://www.asml.com/en
https://ieeexplore.ieee.org/document/4529026
https://www.researchgate.net/figure/shows-the-CD-swing-curve-for-60-and-220-nm-alt-PSM-lines-on-137-nm-thick-oxide-wafers-as_fig15_230779234
