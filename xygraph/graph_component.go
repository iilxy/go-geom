package xygraph

import "github.com/twpayne/go-geom"

type GraphComponent interface {
	// Coordinate returns a coordinate in this component (or nil, if there are none)
	Coordinate() geom.Coord
	// computeIM computes the contribution to an IM for this component
	computeIM(im IntersectionMatrix)
	// isIsolated determins if the component is an isolated component.
	// An isolated component is one that does not intersect or touch any other
	// component.  This is the case if the label has valid locations for
	// only a single Geometry.
	isIsolated() bool
}
type CommonGraphComponent struct {
	label                                          Label
	isInResult, isCovered, isCoveredSet, isVisited bool
}
