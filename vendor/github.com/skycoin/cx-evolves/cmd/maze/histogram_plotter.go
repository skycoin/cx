package maze

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	histoSaveDirectory = "./histogram_data/"
)

func histogramPlot(values plotter.Values, title, saveLocation string) {
	p := plot.New()
	p.Title.Text = title

	hist, err := plotter.NewHist(values, 500)
	if err != nil {
		panic(err)
	}
	p.Add(hist)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, saveLocation); err != nil {
		panic(err)
	}
}
