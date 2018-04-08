package graph

import (
	"image"
	"image/png"
	"io"

	"github.com/Jacobious52/staticmbot/pkg/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

const dpi = 96

type Scatter struct {
	stats model.Timeseries
}

func (s *Scatter) Render(w io.Writer) error {
	exporter := model.UniformUserExporter{Stats: s.stats}
	xVals, yVals := exporter.Export()

	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = "Message Frequency"
	p.X.Label.Text = "Days"
	p.Y.Label.Text = "Messages"

	var series []interface{}
	for user, vals := range yVals {
		xys := make(plotter.XYs, len(xVals))
		for i := range xVals {
			xys[i].X, xys[i].Y = xVals[i], vals[i]
		}
		series = append(series, user, xys)
	}
	err = plotutil.AddLinePoints(p, series...)
	if err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, 16*dpi, 9*dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	err = png.Encode(w, img)
	return err
}
