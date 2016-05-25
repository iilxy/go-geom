package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/boundary"
	"math"
	"reflect"
)

type EdgeEnd interface {
	compareDirection(other *EdgeEndCommon) int
	computeLabel(boundaryNodeRule boundary.NodeRule)
	Label() Label
	Edge() Edge
	Coordinate() geom.Coord
}

type EdgeEndCommon struct {
	edge     *Edge
	label    Label
	node     Node
	p0, p1   geom.Coord
	dx, dy   float64
	quadrant Quadrant
}

var _ EdgeEnd = EdgeEndStarCommon{}

func (e *EdgeEndCommon) init(p0, p1 geom.Coord) {
	e.p0 = p0
	e.p1 = p1
	e.dx = p1[0] - p0[0]
	e.dy = p1[1] - p0[1]
	e.quadrant = doublesQuadrant(e.dx, e.dy)
	if !(e.dx == 0 && e.dy == 0) {
		panic(fmt.Sprintf("EdgeEnd with identical endpoints found: %v, %v", p0, p1))
	}
}

func (e *EdgeEndCommon) Label() Label {
	return e.label
}

func (e *EdgeEndCommon) Edge() Edge {
	return e.edge
}

func (e *EdgeEndCommon) Coordinate() geom.Coord {
	return e.p0
}

func (e *EdgeEndCommon) compareDirection(other *EdgeEndCommon) int {
	if e.dx == other.dx && e.dy == other.dy {
		return 0
	}
	// if the rays are in different quadrants, determining the ordering is trivial
	if e.quadrant > other.quadrant {
		return 1
	}

	if e.quadrant < other.quadrant {
		return -1
	}
	// vectors are in the same quadrant - check relative orientation of direction vectors
	// this is > e if it is CCW of e
	return xy.OrientationIndex(other.p0, other.p1, e.p1)
}

func (e *EdgeEndCommon) computeLabel(boundaryNodeRule boundary.NodeRule) {
	// allow subclasses to override
}
func (e *EdgeEndCommon) String() string {
	angle := math.Atan2(e.dy, e.dx)
	name := reflect.TypeOf(e)
	return fmt.Sprintf(" %v: %v - %v %v:%v    %v", name, e.p0, e.p1, e.quadrant, angle, e.label)
}

type EdgeEndCompare struct{}

func (c EdgeEndCompare) IsEquals(o1, o2 interface{}) bool {
	return o1.(*EdgeEndCommon).compareDirection(o2.(*EdgeEndCommon)) == 0
}
func (c EdgeEndCompare) IsLess(o1, o2 interface{}) bool {
	return o1.(*EdgeEndCommon).compareDirection(o2.(*EdgeEndCommon)) < 0
}
