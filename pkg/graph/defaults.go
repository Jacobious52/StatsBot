package graph

import "github.com/Jacobious52/staticmbot/pkg/model"

// Default graphers
var Default = map[string]func(model.Timeseries) Grapher{
	"alphav1scatter": MakeAlphaV1Scatter,
	"old":            MakeAlphaV1Scatter,
	"v1":             MakeAlphaV1Scatter,
	"scatter":        MakeScatter,
	"new":            MakeScatter,
	"v2":             MakeScatter,
}

func MakeAlphaV1Scatter(s model.Timeseries) Grapher {
	return &AlphaV1Scatter{s}
}

func MakeScatter(s model.Timeseries) Grapher {
	return &Scatter{s}
}
