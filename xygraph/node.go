package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/location"
)

type Node struct {
	commonGraphComponent
	coord geom.Coord
	edges edgeEndStar
}

var _ GraphComponent = &Node{}

func (n *Node) getCoord() geom.Coord {
	return n.coord
}

func (n *Node) getEdges() edgeEndStar {
	if n.edges == nil {
		n.edges = newEdgeEndStarCommon()
	}

	return n.edges
}
// isIncidentEdgeInResult Tests whether any incident edge is flagged as being in the result.
//
// This test can be used to determine if the node is in the result,
// since if any incident edge is in the result, the node must be in the result as well.
func (n *Node) isIncidentEdgeInResult() bool {
	result := false

	n.getEdges().iterate(func(e edgeEnd) bool {
		de := e.(*directedEdge)
		if de.edge.isInResult {
			result = true
			return false
		}
		return true
	})

	return result
}

func (n *Node) isIsolated() bool {
	return n.label.getGeometryCount() == 1
}

func (n *Node) computeIM(im IntersectionMatrix) {
	// do nothing
}

// add sdds the edge to the list of edges at this node
func (n *Node) add(e edgeEnd) {
	// Assert: start pt of e is equal to node point
	n.getEdges().insert(e)
	e.setNode(n)
}

func (n *Node) mergeNodeLabels(other Node) {
	n.mergeLabel(other.label)
}

func (n *Node) mergeLabel(label2 *Label) {
	for i := 0; i < 2; i++ {
		loc := n.computeMergedLocation(label2, i)
		thisLoc := n.label[i][ON]
		if thisLoc == location.None {
			n.label[i] = newOnTopologyLocation(loc)
		}
	}
}

func (n *Node) setLabel(argIndex int, onLocation location.Type) {
	if n.label == nil {
		n.label = NewHomogeneousLabel(onLocation)
	} else {
		n.label[argIndex] = newOnTopologyLocation(onLocation)
	}
}

// setLabelBoundary updates the label of a node to BOUNDARY, obeying the mod-2 boundaryDetermination rule.
func (n *Node) setLabelBoundary(argIndex int) {
	if n.label == nil {
		return
	}

	// determine the current location for the point (if any)
	loc := location.None
	if n.label != nil {
		loc = n.label[argIndex][ON]
	}
	// flip the loc
	var newLoc location.Type
	switch loc {
	case location.Boundary:
		newLoc = location.Interior
	case location.Interior:
		newLoc = location.Boundary
	default:
		newLoc = location.Boundary
	}
	n.label[argIndex] = newOnTopologyLocation(newLoc)
}

func (n *Node) computeMergedLocation(label2 *Label, eltIndex int) location.Type {
	loc := location.None
	loc = n.label[eltIndex][ON]
	if !label2[eltIndex].isNull() {
		nLoc := label2[eltIndex][ON]
		if loc != location.Boundary {
			loc = nLoc
		}
	}
	return loc
}

func (n *Node) String() string {
	return fmt.Sprintf("node %v lbl: %v", n.coord, n.label)
}
