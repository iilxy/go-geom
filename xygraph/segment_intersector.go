package xygraph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/lineintersection"
	"math"
)

type SegmentIntersector struct {
	hasIntersection, hasProper, hasProperInterior,
	includeProper, recordIsolated, isSelfIntersection bool
	properIntersectionPoint geom.Coord
	lineIntersection        lineintersection.Result

	// Boundary nodes
	bdyNodes [2][]*Node
}

func (si *SegmentIntersector) isTrivialIntersection(e0 *Edge, segIndex0 int, e1 *Edge, segIndex1 int) bool {
	if e0 == e1 {
		if len(si.lineIntersection.Intersection()) == 1 {
			if isAdjacentSegments(segIndex0, segIndex1) {
				return true
			}
			if e0.isClosed() {
				maxSegIndex := len(e0.pts) - 1
				if (segIndex0 == 0 && segIndex1 == maxSegIndex) || (segIndex1 == 0 && segIndex0 == maxSegIndex) {
					return true
				}
			}
		}
	}
	return false
}

func isAdjacentSegments(i1, i2 int) bool {
	return math.Abs(float64(i1-i2)) == 1
}

func (si *SegmentIntersector) addIntersections(e0 *Edge, segIndex0 int, e1 *Edge, segIndex1 int) {
	if e0 == e1 && segIndex0 == segIndex1 {
		return
	}

	p00 := e0.pts[segIndex0]
	p01 := e0.pts[segIndex0+1]
	p10 := e1.pts[segIndex1]
	p11 := e1.pts[segIndex1+1]

	si.lineIntersection = xy.LinesIntersection(p00, p01, p10, p11)

	//if (li.hasIntersection() && li.isProper()) Debug.println(li);
	/**
	 *  Always record any non-proper intersections.
	 *  If includeProper is true, record any proper intersections as well.
	 */
	if si.lineIntersection.HasIntersection() {
		if si.recordIsolated {
			e0.isolated = false
			e1.isolated = false
		}

		// if the segments are adjacent they have at least one trivial intersection,
		// the shared endpoint.  Don't bother adding it if it is the
		// only intersection.
		if !si.isTrivialIntersection(e0, segIndex0, e1, segIndex1) {
			si.hasIntersection = true
			if si.includeProper || !si.lineIntersection.IsProper() {
				//Debug.println(li);
				e0.addIntersections(si.lineIntersection, segIndex0, 0)
				e1.addIntersections(si.lineIntersection, segIndex1, 1)
			}
			if si.lineIntersection.IsProper() {
				si.properIntersectionPoint = si.lineIntersection.Intersection()[0]
				si.hasProper = true
				if !isBoundaryPoint(si.lineIntersection, si.bdyNodes) {
					si.hasProperInterior = true
				}
			}
		}
	}
}

func isBoundaryPoint(li lineintersection.Result, bdyNodes [2][]*Node) bool {
	if bdyNodes[0] == nil || bdyNodes[1] == nil {
		return false
	}

	if pointIsOnBoundary(li, bdyNodes[0]) {
		return true
	}
	if pointIsOnBoundary(li, bdyNodes[1]) {
		return true
	}
	return false
}

func pointIsOnBoundary(li lineintersection.Result, bdyNodes []*Node) bool {
	for _, node := range bdyNodes {
		if li.IsIntersectionPoint(node.coord) {
			return true
		}
	}
	return false
}
