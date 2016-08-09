package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom/xy/location"
)

type directedEdge struct {
	edgeEndCommon

	isForward, isInResult, isVisited bool
	sym,
	// the next edge in the edge ring for the polygon containing this edge
	next,
	// the next edge in the MinimalEdgeRing that contains this edge
	nextMin *directedEdge
	// the EdgeRing that this edge is part of
	edgeRing,
	// the MinimalEdgeRing that this edge is part of
	minEdgeRing *edgeRing
	// The depth of each side (position) of this edge.
	//  The 0 element of the array is never used.
	depth [3]int
}

var _ edgeEnd = &directedEdge{}

func newDirectedEdge(edge *Edge, isForward bool) *directedEdge {
	de := &directedEdge{
		edgeEndCommon: edgeEndCommon{
			edge: edge,
		},
		depth:     [3]int{0, -999, -999},
		isForward: isForward,
	}

	stride := edge.layout.Stride()
	doubleStride := stride*2
	if isForward {
		de.init(edge.pts[:stride], edge.pts[stride:doubleStride])
	} else {
		idx := len(edge.pts) - stride
		de.init(edge.pts[idx:], edge.pts[idx -stride:idx])
	}
	de.computeDirectedLabel()

	return de
}

// Computes the factor for the change in depth when moving from one location to another.
// E.g. if crossing from the INTERIOR to the EXTERIOR the depth decreases, so the factor is -1
func depthFactor(currLocation, nextLocation location.Type) int {
	switch {
	case currLocation == location.Exterior && nextLocation == location.Interior:
		return 1
	case currLocation == location.Interior && nextLocation == location.Exterior:
		return -1
	default:
		return 0
	}
}

func (de *directedEdge) setDepth(pos position, depthVal int) {
	if de.depth[pos] != -999 {
		if de.depth[pos] != depthVal {
			panic(fmt.Sprintf("assigned depths do not match: %v", de.p0))
		}
	}
	de.depth[pos] = depthVal
}

func (de *directedEdge) getDepthDelta() int {
	depthDelta := de.edge.depthDelta
	if !de.isForward {
		depthDelta = -depthDelta
	}
	return depthDelta
}

// setVisitedEdge marks both DirectedEdges attached to a given Edge.
// This is used for edges corresponding to lines, which will only appear oriented in a single direction in the result.
func (de *directedEdge) setVisitedEdge(isVisited bool) {
	de.isVisited = isVisited
	de.sym.isVisited = isVisited
}

// isLineEdge determines if this edige is a line edge.  This edge is a line edge if
//  * at least one of the labels is a line label
//  * any labels which are not line labels have all Locations = EXTERIOR
func (de *directedEdge) isLineEdge() bool {
	isLine := de.label[0].isLine() || de.label[1].isLine()
	isExteriorIfArea0 := !de.label[0].isArea() || de.label[0].allPositionsEqual(location.Exterior)
	isExteriorIfArea1 := !de.label[1].isArea() || de.label[1].allPositionsEqual(location.Exterior)

	return isLine && isExteriorIfArea0 && isExteriorIfArea1
}

// This is an interior Area edge if
//  * its label is an Area label for both Geometries
//  * and for each Geometry both sides are in the interior.
func (de *directedEdge) isInteriorAreaEdge() bool {
	isInteriorAreaEdge := true
	for i := 0; i < 2; i++ {
		if !(de.label[i].isArea() && de.label[i][LEFT] == location.Interior && de.label[i][RIGHT] == location.Interior) {
			isInteriorAreaEdge = false
		}
	}
	return isInteriorAreaEdge
}

// Compute the label in the appropriate orientation for this DirEdge
func (de *directedEdge) computeDirectedLabel() {
	if de.edge.label != nil {
		de.label = NewLabelFromTemplate(de.edge.label)

		if !de.isForward {
			de.label.flip()
		}
	}
}

// setEdgeDepths sets both edge depths.  One depth for a given side is provided.  The other is
// computed depending on the Location transition and the depthDelta of the edge.
func (de *directedEdge) setEdgeDepths(pos position, depth int) {
	// get the depth transition delta from R to L for this directed Edge
	depthDelta := de.edge.depthDelta
	if !de.isForward {
		depthDelta = -depthDelta
	}

	// if moving from L to R instead of R to L must change sign of delta
	directionFactor := 1
	if pos == LEFT {
		directionFactor = -1
	}

	oppositePos := pos.opposite()
	delta := depthDelta * directionFactor

	oppositeDepth := depth + delta
	de.setDepth(pos, depth)
	de.setDepth(oppositePos, oppositeDepth)
}

func (de *directedEdge) String() string {
	s := fmt.Sprintf("%v %v/%v (%v) ", de.edgeEndCommon.String(), de.depth[LEFT], de.depth[RIGHT], de.getDepthDelta())

	if de.isInResult {
		s = s + " inResult"
	}
	return s
}
