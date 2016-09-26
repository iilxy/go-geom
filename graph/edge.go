package graph

import "github.com/twpayne/go-geom"

// Edge represents an edge in a graph
type Edge struct {
	GraphableImpl
}

var _ Graphable = Edge{}

func (e Edge) Intersections() []EdgeIntersection {

}

type EdgeIntersection struct {
	// SegmentIndex is the index of the Edge LineString on which the intersection lies
	SegmentIndex int
	// Coord is the location of the intersection
	Coord geom.Coord
	// Dist is the distance along the line segment (indexed by SegmentIndex) of the intersection (Coord).
	Dist float64
	// The Edge the intersectino lies on
	Edge geom.LineString
}

type edgeIntersections struct {
}
