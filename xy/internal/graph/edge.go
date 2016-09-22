package graph

import (
	"bytes"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/lineintersection"
)

type Edge struct {
	BasicComponent
	layout   geom.Layout
	pts      []float64
	env      *geom.Bounds
	eiList   edgeIntersectionList
	name     string
	mce      *monotoneChainEdge
	isolated bool
	depth    Depth
	// depthDelta is the change in depth as an edge is crossed from R to L
	depthDelta int
}

var _ Component = &Edge{}

func NewEdge(layout geom.Layout, pts []float64, label *Label) *Edge {
	return &Edge{
		BasicComponent: BasicComponent{
			label: label,
		},
		layout:   layout,
		pts:      pts,
		isolated: true,
		depth:    newDepth(),
	}
}

func (e *Edge) bounds() *geom.Bounds {
	// compute envelope lazily
	if e.env == nil {
		e.env = geom.NewBounds(geom.XY)
		stride := e.layout.Stride()
		for i := 0; i < len(e.pts)-stride; i += stride {
			c := geom.Coord(e.pts[i : i+stride])
			e.env.ExtendWithCoord(c)
		}
	}
	return e.env
}

func (e *Edge) Coord(idx int) geom.Coord {
	return geom.Coord(e.pts[idx : idx+e.layout.Stride()])
}
func (e *Edge) NumCoords() int {
	return len(e.pts) / e.layout.Stride()
}
func (e *Edge) MonotoneChainEdge() *monotoneChainEdge {
	if e.mce == nil {
		e.mce = newMonotoneChainEdge(e)
	}

	return e.mce
}

func (e *Edge) isClosed() bool {
	return xy.Equal(e.pts, 0, e.pts, len(e.pts)-e.layout.Stride()-1)
}

// isCollapsed returns true if the edge is an Area edge and it consists of two
// segments which are equal and opposite (eg a zero-width V).
func (e *Edge) isCollapsed() bool {
	if !e.label.isArea() {
		return false
	}
	if len(e.pts) != 3 {
		return false
	}
	if xy.Equal(e.pts, 0, e.pts, 2*e.layout.Stride()) {
		return true
	}
	return false
}

func (e *Edge) FirstCoord() geom.Coord {
	if len(e.pts) > 0 {
		return geom.Coord(e.pts[:e.layout.Stride()])
	}
	return nil
}
func (e *Edge) isIsolated() bool {
	return e.isolated
}
func (e *Edge) collapsedEdge() *Edge {
	newPts := make([]float64, e.layout.Stride()*2)
	copy(newPts, e.pts)

	return NewEdge(e.layout, newPts, e.label.toLineLabel())
}

// Adds EdgeIntersections for one or both
// intersections found for a segment of an edge to the edge intersection list.
func (e *Edge) addIntersections(li lineintersection.Result, segmentIndex, geomIndex int) {
	intersections := li.Intersection()
	for i := 0; i < len(intersections); i++ {
		e.addIntersection(intersections[i], segmentIndex, geomIndex)
	}
}

func (e *Edge) addIntersection(intPt geom.Coord, segmentIndex, geomIndex int) {
	normalizedSegmentIndex := segmentIndex

	stride := e.layout.Stride()
	segmentStart := geom.Coord(e.pts[geomIndex : geomIndex+stride])
	segmentEnd := geom.Coord(e.pts[geomIndex+stride : geomIndex+stride+stride])
	dist := xy.DistanceFromPointToLine(intPt, segmentStart, segmentEnd)

	// normalize the intersection point location
	nextSegIndex := normalizedSegmentIndex + 1
	if nextSegIndex < len(e.pts) {
		arrayIdx := nextSegIndex * stride
		nextPt := e.pts[arrayIdx : arrayIdx+stride]

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

func updateIM(label *Label, im IntersectionMatrix) {
	im.SetAtLeastIfValid(int(label[0][OnLabel]), int(label[1][OnLabel]), 1)
	if label.isArea() {
		im.SetAtLeastIfValid(int(label[0][LeftOfLabel]), int(label[1][LeftOfLabel]), 2)
		im.SetAtLeastIfValid(int(label[0][RightOfLabel]), int(label[1][RightOfLabel]), 2)
	}

}

func (e *Edge) isPointwiseEqual(o Edge) bool {
	if len(e.pts) != len(o.pts) {
		return false
	}

	stride := e.layout.Stride()

	for i := 0; i < len(e.pts); i += stride {
		if !xy.Equal(e.pts, 0, o.pts, i*stride) {
			return false
		}
	}
	return true
}

func (e *Edge) equal(o Edge) bool {
	if len(e.pts) != len(o.pts) {
		return false
	}

	stride := e.layout.Stride()

	isEqualForward := true
	isEqualReverse := true
	iRev := len(e.pts)
	for i := 0; i < len(e.pts); i += stride {
		if !xy.Equal(e.pts, i, o.pts, i+stride) {
			isEqualForward = false
		}

		iRev = iRev - stride
		if !xy.Equal(e.pts, i, e.pts, iRev) {
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
	stride := e.layout.Stride()

	for i := 0; i < len(e.pts); i += stride {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%v %v", e.pts[i:i+stride], e.pts[i+stride:i+stride+stride]))
	}
	buf.WriteString(fmt.Sprintf(")  %v %v", e.label, e.depthDelta))
	return buf.String()
}

func (e *Edge) StringReverse() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("edge %v: ", e.name))
	for i := len(e.pts) - 1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%v ", e.pts[i]))
	}
	buf.WriteString("\n")
	return buf.String()
}
