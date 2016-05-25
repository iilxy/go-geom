package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/location"
)

type Node struct {
	CommonGraphComponent
	coord geom.Coord
	edges EdgeEndStar
}

var _ GraphComponent = &Node{}

func (n *Node) getCoordinate() geom.Coord {
	return n.coord
}

// isIncidentEdgeInResult Tests whether any incident edge is flagged as being in the result.
//
// This test can be used to determine if the node is in the result,
// since if any incident edge is in the result, the node must be in the result as well.

func (n *Node) isIncidentEdgeInResult() bool {
	for _, e := range n.edges {
		de := e.(DirectedEdge)
		if de.getEdge().isInResult() {
			return true
		}
	}
	return false
}

func (n *Node) isIsolated() bool {
	return n.label.getGeometryCount() == 1
}

func (n *Node) computeIM(im IntersectionMatrix) {
	// do nothing
}

// add sdds the edge to the list of edges at this node
func (n *Node) add(e EdgeEnd) {
	// Assert: start pt of e is equal to node point
	n.edges.insert(e)
	e.setNode(this)
}

func (n *Node) mergeNodeLabels(n Node) {
	n.mergeLabel(n.label)
}

func (n *Node) mergeLabel(label2 Label) {
	for i := 0; i < 2; i++ {
		loc := n.computeMergedLocation(label2, i)
		thisLoc := n.label.getLocation(i)
		if thisLoc == location.None {
			n.label[i] = NewOnTopologyLocation(loc)
		}
	}
}

func (n *Node) setLabel(argIndex int, onLocation location.Type) {
	if n.label == nil {
		n.label = NewHomogeneousLabel(onLocation)
	} else {
		n.label[argIndex] = NewOnTopologyLocation(onLocation)
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
		loc = n.label[argIndex]
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
	n.label[argIndex] = NewOnTopologyLocation(newLoc)
}

func (n *Node) computeMergedLocation(label2 Label, eltIndex int) {
	loc := location.None
	loc = location.Type(n.label[eltIndex])
	if !label2[eltIndex].isNull() {
		nLoc := location.Type(label2[eltIndex])
		if loc != location.Boundary {
			loc = nLoc
		}
	}
	return loc
}

func (n *Node) String() string {
	return fmt.Sprintf("node %v lbl: %v", n.coord, n.label)
}
