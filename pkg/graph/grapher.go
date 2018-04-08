package graph

import "io"

// Grapher renders a graph to a writer
type Grapher interface {
	Render(io.Writer) error
}
