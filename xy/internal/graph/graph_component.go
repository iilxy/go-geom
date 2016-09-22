package graph

import "github.com/twpayne/go-geom"

type Component interface {
	// Coordinate returns a coordinate in this component (or nil, if there are none)
	FirstCoord() geom.Coord
	// ItersectionMatrix computes the contribution to an IM for this component
	computeIM(im IntersectionMatrix)
	// isIsolated determins if the component is an isolated component.
	// An isolated component is one that does not intersect or touch any other
	// component.  This is the case if the label has valid locations for
	// only a single Geometry.
	isIsolated() bool
}
type BasicComponent struct {
	label                                          *Label
	isInResult, isCovered, isCoveredSet, isVisited bool
}
