package noding

import "github.com/twpayne/go-geom"

type SegmentString interface {
	// Gets the user-defined data for this segment string.
	GetData() interface{}
	// Sets the user-defined data for this segment string.
	SetData(data interface{})
	Size() int
	GetCoordinates() []geom.Coord
	IsClosed() bool
}

type NodeableSegmentString interface {
	SegmentString
	AddIntersection(intPt geom.Coord, segmentIndex int)
}

type NodedSegmentString struct {
	nodeList SegmentNodeList
	pts      []geom.Coord
	data     interface{}
}
