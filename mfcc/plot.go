package mfcc

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func PlotMFCC(inputSignal [][]float64) {
	p := plot.New()

	p.Title.Text = "FFT Signal Visualization"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	length := len(inputSignal)

	err := plotutil.AddLinePoints(p,
		"First", setPoints(length, inputSignal))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "plotMFCC.png"); err != nil {
		panic(err)
	}
}

func setPoints(length int, inputSignal [][]float64) plotter.XYs {
	pts := make(plotter.XYs, length)
	maxVal := 0.0
	for i := range inputSignal {
		for j := 0; j < len(inputSignal[i]); j++ {
			if inputSignal[i][j] > maxVal {
				maxVal = inputSignal[i][j]
			}

		}
		pts[i].X = float64(i)
		pts[i].Y = maxVal
		maxVal = 0.0
	}

	return pts
}
