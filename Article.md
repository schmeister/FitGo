# Solving a Swing Curve with GO and go-hep and GoNum

My first professional task as a Software Engineer came just over 30 years ago. I was in college working towards my BS in Computer science and had completed Physics I & II, Calculus I & II, and a number of other Computer Science pre-requisites, and had recently been promoted to Process Technician in a Semi-conductor company. 

Going back even farther, I had been programming for quite some time. In the early 80's while still in middle school my dad purchased our first computer, an Atari 400 8-bit computer with 4k RAM. I pleaded and eventually got a few upgrades: the "Basic" cartridge, a tape drive, and a subscription to Antic magazine. I had typed in many of the programs, but the one I remember the most was Bats (Antic, Stan Ockers - V1 No 5) and modified that code dozens of times to better understand what was happening and tailor the game play to what I wanted.

Returning to the "semi-conductor days", I worked in various operator roles and departments - from packaging (installing microprocessors into the package, ready for installation in circuit boards) to testing. The latest step was in the main clean room where the actual laying down of circuits on the silicon wafers occured, I was a Chemical Technician - the person who handled and mixed all the dangerous chemicals while I wore many layers of protective gear and spent a lot of time with my hands in high ventilation hoods. Then my big break came as a Process Technician in the Photolithography department, my RAA (Responsibility, Accountability, and Authority) took a big step up. Though many details have been lost over time, I will try to resurrect some of what I remember but with a modern flare and tool suite.

We were staffed such that there were one or two Process technicians per shift in each part of the FAB (Fabrication Plant), my shift was 3rd shift (11pm to 7am) as the sole process technician. As part of my RAA, I needed to keep the Photolithography line running to the best of my ability and handling common deviations. One of those standard engineering processes was when the resist packages would run out and need to be replaced. For the uninformed, a resist is a photosensitive liquid that is evenly dispersed on the surface of a silicon wafer. The thickness of the resist needs to be very accurate and there are many engineering processes necessary to ensure the proper thickness. Ensuring this thickness is where my Software Engineering skill was brought forth. Even if the resist is from the same batch from the manufacturer, there are always some variations and the development line needs to be validated before any further products could be processed.

Thus came the Swing Curve: A sine wave graph that is fit from a sampling of different resist thicknesses on reference wafers. 
**Show a swing curve with some offset and decay**


The resist is dispensed on a number wafers with varying spin speeds that should obtain the approximate required thickness of the resist with a known viscosity (lower viscosity implies a slower spin speed than thicker viscosity to obtain the same thickness). The thickness was measured with highly accurate measuring tools (https://www.kla.com/), exposed (https://www.asml.com/en), developed (DNS), and finally verified with an SEM (Scanning Electron Microscope). The intermediate step would look something like this:

**Insert raw data graph**

At this point we could probably come up with an approximation of the optimal spin speed, but that is not good enough, we need as close as possible! With that said, the next task is to use those points to find the **exact** best speed for our desired thickness.

Lets take a step back and wonder why this is so important. There are a signficant number of physics at play, but we will focus on one of the two most important aspects of laying down an integrated circuit (https://ieeexplore.ieee.org/document/4529026): resist thickness. (https://en.wikipedia.org/wiki/Thin-film_interference)

From the graph above, what should we choose as the optimal spin speed to get the best thickness of the resist? Well first off, we want the most stable location, meaning that any variations in either spin speed or viscosity (potentially due to variation in temperature) result in the least amount of change. To do that we need to choose where the slope of the line is near zero. This tell us that we need two things now. A common formula that we can use to model and fit those points, then the derivative of that formula (remember that a derivative of a formula will help you determine the slope of a line.)

In the first case, we will use this formula to model a sine curve: f(x) = Amplitude * sin(Lambda * x + Shift) * e^((-x^2)*Decay).
Simple enough, now to find the optimal slope, let's compute the derivative, but we don't need the full derivative, we only need to determine when the sine function has zero slope, and this is the first derviative of f(x): f'(x) = cos(Lambda * x + Shift) = 0. By setting f'(x) to 0, that means we are going to solve f'(x) where f'(x) == 0.

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
