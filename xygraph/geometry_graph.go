package xygraph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/boundary"
)

// GeometryGraph is a graph that models a given Geometry
type GeometryGraph struct {
	parentGeom                   geom.T
	lineEdgeMap                  map[*geom.LineString]Edge
	boundaryNodeRule             boundary.NodeRule
	useBoundaryDeterminationRule bool
	argIndex                     int
	boundaryNodes                []Node
	hasTooFewPoints              bool
	invalidPoint                 geom.Coord
}
