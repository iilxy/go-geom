package xygraph

import (
	"bytes"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/lineintersection"
)

type Edge struct {
	CommonGraphComponent
	pts        []geom.Coord
	env        geom.Bounds
	eiList     edgeIntersectionList
	name       string
	mce        MonotoneChainEdge
	isIsolated bool
	depth      depth
	// depthDelta is the change in depth as an edge is crossed from R to L
	depthDelta int
}

var _ GraphComponent = &Edge{}

func NewEdge(pts []geom.Coord, label Label) *Edge {
	return &Edge{
		CommonGraphComponent{
			label: label,
		},
		pts:        pts,
		isIsolated: true,
		depth:      newDepth(),
		env:        nil,
	}
}

func (e *Edge) Envelope() {
	// compute envelope lazily
	if e.env == nil {
		e.env = geom.NewBounds(geom.XY)
		for _, c := range e.pts {
			e.env.Extend(c)
		}
	}
	return e.env
}

func (e *Edge) MonotoneChainEdge() {
	if e.mce == nil {
		e.mce = newMonotoneChainEdge(e)
	}

	return e.mce
}

func (e *Edge) isClosed() {
	return xy.Equal(e.pts[0], 0, e.pts[len(e.pts-1)], 0)
}

// isCollapsed returns true if the edge is an Area edge and it consists of two
// segments which are equal and opposite (eg a zero-width V).
func (e *Edge) isCollapsed() {
	if !e.label.isArea() {
		return false
	}
	if len(e.pts) != 3 {
		return false
	}
	if xy.Equal(e.pts[0], 0, e.pts[2], 0) {
		return true
	}
	return false
}

func (e *Edge) getCoordinate() geom.Coord {
	if len(e.pts) > 0 {
		return e.pts[0]
	}
	return nil
}
func (e *Edge) collapsedEdge() {
	newPts := [2]geom.Coord{e.pts[0], e.pts[1]}
	return NewEdge(newPts, e.label.toLineLabel())
}

// Adds EdgeIntersections for one or both
// intersections found for a segment of an edge to the edge intersection list.
func (e *Edge) addIntersections(li lineintersection.Result, segmentIndex, geomIndex int) {
	intersections := li.Intersection()
	for i := 0; i < len(intersections); i++ {
		e.addIntersection(li[i], segmentIndex, geomIndex)
	}
}

func (e *Edge) addIntersection(intPt geom.Coord, segmentIndex, geomIndex int) {
	normalizedSegmentIndex := segmentIndex

	dist := xy.DistanceFromPointToLine(intPt, e.pts[geomIndex], e.pts[geomIndex+1])

	// normalize the intersection point location
	nextSegIndex := normalizedSegmentIndex + 1
	if nextSegIndex < len(e.pts) {
		nextPt := e.pts[nextSegIndex]

		// Normalize segment index if intPt falls on vertex
		// The check for point equality is 2D only - Z values are ignored
		if xy.Equal(intPt, 0, nextPt, 0) {
			normalizedSegmentIndex = nextSegIndex
			dist = 0.0
		}
	}
	/**
	* Add the intersection point to edge intersection list.
	 */
	e.eiList.add(intPt, normalizedSegmentIndex, dist)
}

func (e *Edge) computeIM(im IntersectionMatrix) {
	updateIM(e.label, im)
}

func updateIM(label Label, im IntersectionMatrix) {
	im.setAtLeastIfValid(label[0][ON], label[1][ON], 1)
	if label.isArea() {
		im.setAtLeastIfValid(label[0][LEFT], label[1][LEFT], 2)
		im.setAtLeastIfValid(label[0][RIGHT], label[1][RIGHT], 2)
	}

}

func (e *Edge) isPointwiseEqual(o Edge) bool {
	if len(e.pts) != len(o.pts) {
		return false
	}

	for i := 0; i < len(e.pts); i++ {
		if !xy.Equal(e.pts[i], 0, o.pts[i], 0) {
			return false
		}
	}
	return true
}

func (e *Edge) equal(o Edge) bool {
	if len(e.pts) != len(o.pts) {
		return false
	}

	isEqualForward := true
	isEqualReverse := true
	iRev := len(e.pts)
	for i := 0; i < len(e.pts); i++ {
		if !xy.Equal(e.pts[i], 0, o.pts[i], 0) {
			isEqualForward = false
		}
		otherPt := o.pts[iRev]
		iRev = iRev - 1
		if !xy.Equal(e.pts[i], 0, otherPt, 0) {
			isEqualReverse = false
		}
		if !isEqualForward && !isEqualReverse {
			return false
		}
	}

	return true
}

func (e *Edge) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("edge %v:", e.name))
	buf.WriteString("LINESTRING (")
	for i, v := range e.pts {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%v %v", v[0], v[1]))
	}
	buf.WriteString(fmt.Sprintf(")  %v %v", e.label, e.depthDelta))
	return buf.String()
}

func (e *Edge) StringReverse() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("edge %v: ", e.name))
	for i := len(e.pts) - 1; i >= 0; i-- {
		buf.WriteString(e.pts[i] + " ")
		buf.WriteString(" ")
	}
	buf.WriteString("\n")
	return buf.String()
}
