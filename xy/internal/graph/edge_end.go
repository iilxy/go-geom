package graph

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/boundary"
	"math"
	"reflect"
)

type edgeEnd interface {
	compareDirection(other edgeEnd) int
	computeLabel(boundaryNodeRule boundary.NodeRule)
	getLabel() *Label
	getEdge() *Edge
	setNode(node *Node)
	getCoord() geom.Coord
	commonData() *edgeEndCommon
}

type edgeEndCommon struct {
	edge     *Edge
	label    *Label
	node     *Node
	p0, p1   geom.Coord
	dx, dy   float64
	quadrant quadrant
}

var _ edgeEnd = &edgeEndCommon{}

func (e *edgeEndCommon) init(p0, p1 geom.Coord) {
	e.p0 = p0
	e.p1 = p1
	e.dx = p1[0] - p0[0]
	e.dy = p1[1] - p0[1]
	e.quadrant = doublesQuadrant(e.dx, e.dy)
	if e.dx == 0 && e.dy == 0 {
		panic(fmt.Sprintf("EdgeEnd with identical endpoints found: %v, %v", p0, p1))
	}
}

func (e *edgeEndCommon) getLabel() *Label {
	return e.label
}

func (e *edgeEndCommon) setNode(node *Node) {
	e.node = node
}

func (e *edgeEndCommon) getEdge() *Edge {
	return e.edge
}

func (e *edgeEndCommon) commonData() *edgeEndCommon {
	return e
}

func (e *edgeEndCommon) getCoord() geom.Coord {
	return e.p0
}

func (e *edgeEndCommon) compareDirection(other edgeEnd) int {
	commonData := other.commonData()
	if e.dx == commonData.dx && e.dy == commonData.dy {
		return 0
	}
	// if the rays are in different quadrants, determining the ordering is trivial
	if e.quadrant > commonData.quadrant {
		return 1
	}

	if e.quadrant < commonData.quadrant {
		return -1
	}
	// vectors are in the same quadrant - check relative orientation of direction vectors
	// this is > e if it is CCW of e
	return int(xy.OrientationIndex(commonData.p0, commonData.p1, e.p1))
}

func (e *edgeEndCommon) computeLabel(boundaryNodeRule boundary.NodeRule) {
	// allow subclasses to override
}
func (e *edgeEndCommon) String() string {
	angle := math.Atan2(e.dy, e.dx)
	name := reflect.TypeOf(e)
	return fmt.Sprintf(" %v: %v - %v %v:%v    %v", name, e.p0, e.p1, e.quadrant, angle, e.label)
}

type EdgeEndCompare struct{}

func (c EdgeEndCompare) IsEquals(o1, o2 interface{}) bool {
	if o1 == nil && o2 == nil {
		return true
	}
	if (o1 == nil && o2 != nil) || (o1 != nil && o2 == nil) {
		return false
	}
	return o1.(edgeEnd).compareDirection(o2.(edgeEnd)) == 0
}
func (c EdgeEndCompare) IsLess(o1, o2 interface{}) bool {
	return o1 == nil || o1.(edgeEnd).compareDirection(o2.(edgeEnd)) < 0
}
