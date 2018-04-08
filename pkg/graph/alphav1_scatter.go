package graph

import (
	"fmt"
	"io"
	"time"

	"github.com/Jacobious52/staticmbot/pkg/model"
	cg "github.com/wcharczuk/go-chart"
)

type AlphaV1Scatter struct {
	stats model.Timeseries
}

func (a *AlphaV1Scatter) Render(w io.Writer) error {
	exporter := model.UniformUserExporter{a.stats}
	xVals, yVals := exporter.Export()

	var series []cg.Series
	var i int
	for user, vals := range yVals {
		series = append(series, cg.ContinuousSeries{
			Style: cg.Style{
				Show:        true,
				StrokeColor: cg.GetDefaultColor(i).WithAlpha(64),
				FillColor:   cg.GetDefaultColor(i).WithAlpha(64),
			},
			Name:    user,
			XValues: xVals,
			YValues: vals,
		})
		i++
	}

	graph := cg.Chart{
		XAxis: cg.XAxis{
			Name:      "Day of Year",
			NameStyle: cg.StyleShow(),
			Style:     cg.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if f, isFloat := v.(float64); isFloat {
					now := time.Now()
					return time.Date(now.Year(), 1, int(f), 0, 0, 0, 0, time.UTC).Format("Jan 2")
				}
				return fmt.Sprint(v)
			},
		},
		YAxis: cg.YAxis{
			Name:      "Message Count",
			NameStyle: cg.StyleShow(),
			Style:     cg.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if f, isFloat := v.(float64); isFloat {
					return fmt.Sprint(int(f))
				}
				return fmt.Sprint(v)
			},
		},
		Series: series,
		Background: cg.Style{
			Padding: cg.Box{
				Top:  20,
				Left: 20,
			},
		},
	}

	graph.Elements = []cg.Renderable{
		cg.Legend(&graph),
	}

	err := graph.Render(cg.PNG, w)
	return err
}
