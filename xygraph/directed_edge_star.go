package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom/xy/location"
)

const (
	SCANNING_FOR_INCOMING = iota + 1
	LINKING_TO_OUTGOING
)

type DirectedEdgeStar struct {
	EdgeEndStarCommon
	resultAreaEdgeList []DirectedEdge
	label              Label
}

func NewDirectedEdgeStart(edgeList []EdgeEnd) *DirectedEdgeStar {
	star := &DirectedEdgeStar{
		NewEdgeEndStarCommon(edgeList),
		resultAreaEdgeList: nil,
	}

	for _, e := range edgeList {
		star.insert(e)
	}
	return star
}

func (des *DirectedEdgeStar) insert(ee EdgeEnd) {
	de := ee.(DirectedEdge)
	des.insertEdgeEnd(de, de)
}

func (des *DirectedEdgeStar) getOutgoingDegreeInResult() int {
	return des.getOutgoingDegree(func(de DirectedEdge) {
		return de.isInResult
	})
}

func (des *DirectedEdgeStar) getOutgoingDegreeInRing(er EdgeRing) int {
	return des.getOutgoingDegree(func(de DirectedEdge) {
		return de.edgeRing == er
	})
}

func (des *DirectedEdgeStar) getOutgoingDegree(cmp func(DirectedEdge) bool) int {
	degree := 0
	des.edgeMap.Walk(func(key, e interface{}) {
		de := e.(DirectedEdge)
		if cmp(de) {
			degree++
		}
	})

	return degree
}

func (des *DirectedEdgeStar) getRightmostEdge() DirectedEdge {
	edges := des.edgeList
	size := len(edges)
	if size < 1 {
		return nil
	}
	de0 := edges[0].(DirectedEdge)
	if size == 1 {
		return de0
	}
	deLast := edges[size-1].(DirectedEdge)

	quad0 := de0.quadrant
	quad1 := deLast.quadrant
	if Quadrant.isNorthern(quad0) && Quadrant.isNorthern(quad1) {
		return de0
	} else if !Quadrant.isNorthern(quad0) && !Quadrant.isNorthern(quad1) {
		return deLast
	} else {

		if de0.dy != 0 {
			return de0
		} else if deLast.dy != 0 {
			return deLast
		}
	}

	// should never reach here
	panic("Found two horizontal edges incident on node, This should not be possible")
}

// computeLabellingCompute computes the labelling for all dirEdges in this star, as well as the overall labelling
func (des *DirectedEdgeStar) computeLabelling(geom []GeometryGraph) {
	des.EdgeEndStarCommon.computeLabelling(geom)

	// determine the overall labelling for this DirectedEdgeStar
	// (i.e. for the node it is based at)
	des.label = NewLabel(location.None)

	des.edgeMap.Walk(func(key, value interface{}) {
		ee := value.(DirectedEdge)
		e := ee.edge
		eLabel := e.label
		for i := 0; i < 2; i++ {
			eLoc := eLabel[i]
			if eLoc == location.Interior || eLoc == location.Boundary {
				des.label[i] = location.Interior
			}
		}
	})

}

// mergeSymLabels merges the label from the sym dirEdge into the label for each dirEdge in the star
func (des *DirectedEdgeStar) mergeSymLabels() {

	des.edgeMap.Walk(func(key, ee interface{}) {
		de := ee.(DirectedEdge)
		label := de.Label()
		label.merge(de.sym.label)
	})
}

// updateLabelling updates incomplete dirEdge labels from the labelling for the node
func (des *DirectedEdgeStar) updateLabelling(nodeLabel Label) {
	des.edgeMap.Walk(func(key, ee interface{}) {
		de := ee.(DirectedEdge)
		label := de.label
		label[0].setAllLocationsIfNull(nodeLabel[0])
		label[0].setAllLocationsIfNull(nodeLabel[1])
	})
}

func (des *DirectedEdgeStar) getResultAreaEdges() []DirectedEdge {
	if des.resultAreaEdgeList != nil {
		return des.resultAreaEdgeList
	}
	des.resultAreaEdgeList = []DirectedEdge{}

	des.edgeMap.Walk(func(key, ee interface{}) {
		de := ee.(DirectedEdge)
		if de.isInResult() || de.sym.isInResult {
			des.resultAreaEdgeList = append(des.resultAreaEdgeList, de)
		}
	})
	return des.resultAreaEdgeList
}

// Traverse the star of DirectedEdges, linking the included edges together.
// To link two dirEdges, the <next> pointer for an incoming dirEdge
// is set to the next outgoing edge.
//
// DirEdges are only linked if:
// * they belong to an area (i.e. they have sides)
// * they are marked as being in the result
//
// Edges are linked in CCW order (the order they are stored).
// This means that rings have their face on the Right
// (in other words,
// the topological location of the face is given by the RHS label of the DirectedEdge)
//
// PRECONDITION: No pair of dirEdges are both marked as being in the result
func (des *DirectedEdgeStar) linkResultDirectedEdges() error {
	// make sure edges are copied to resultAreaEdges list
	des.getResultAreaEdges()
	// find first area edge (if any) to start linking at
	var firstOut, incoming DirectedEdge = nil, nil
	state := SCANNING_FOR_INCOMING
	// link edges in CCW order
	for i := 0; i < len(des.resultAreaEdgeList); i++ {
		nextOut := des.resultAreaEdgeList[i]
		nextIn := nextOut.sym

		// skip de's that we're not interested in
		if !nextOut.label.isArea() {
			continue
		}

		// record first outgoing edge, in order to link the last incoming edge
		if firstOut == nil && nextOut.isInResult {
			firstOut = nextOut
		}

		switch state {
		case SCANNING_FOR_INCOMING:
			if !nextIn.isInResult {
				continue
			}
			incoming = nextIn
			state = LINKING_TO_OUTGOING
		case LINKING_TO_OUTGOING:
			if !nextOut.isInResult {
				continue
			}
			incoming.next = nextOut
			state = SCANNING_FOR_INCOMING
		}
	}

	if state == LINKING_TO_OUTGOING {
		if firstOut == nil {
			return fmt.Errorf("no outgoing dirEdge found %v", des.Coordinate())
		}
		if !firstOut.isInResult {
			return fmt.Errorf("unable to link last incoming dirEdge")
		}
		incoming.next = firstOut
	}

	return nil
}

func (des *DirectedEdgeStar) linkMinimalDirectedEdges(er EdgeRing) error {

	des.getResultAreaEdges()
	// find first area edge (if any) to start linking at
	var firstOut, incoming DirectedEdge = nil, nil

	state := SCANNING_FOR_INCOMING
	// link edges in CW order
	for i := 0; i < len(des.resultAreaEdgeList); i++ {
		nextOut := des.resultAreaEdgeList[i]
		nextIn := nextOut.sym

		// record first outgoing edge, in order to link the last incoming edge
		if firstOut == nil && nextOut.edgeRing == er {
			firstOut = nextOut
		}

		switch state {
		case SCANNING_FOR_INCOMING:
			if nextIn.edgeRing != er {
				continue
			}
			incoming = nextIn
			state = LINKING_TO_OUTGOING
		case LINKING_TO_OUTGOING:
			if nextOut.edgeRing != er {
				continue
			}
			incoming.nextMin = nextOut
			state = SCANNING_FOR_INCOMING
		}
	}

	if state == LINKING_TO_OUTGOING {
		if firstOut == nil {
			return fmt.Errorf("Did not find an edge for the first outgoing dirEdge")
		}
		if firstOut.edgeRing != er {
			return fmt.Errorf("unable to link last incoming dirEdge")
		}
		incoming.nextMin = firstOut
	}
	return nil
}

func (des *DirectedEdgeStar) linkAllDirectedEdges() {
	// find first area edge (if any) to start linking at
	var prevOut, firstIn DirectedEdge = nil, nil

	// link edges in CW order
	for i := 0; i < len(des.edgeList); i++ {
		nextOut := des.edgeList[i].(DirectedEdge)
		nextIn := nextOut.sym

		if firstIn == nil {
			firstIn = nextIn
		}
		if prevOut != nil {
			nextIn.next = prevOut
		}
		// record outgoing edge, in order to link the last incoming edge
		prevOut = nextOut
	}
	firstIn.next = prevOut
}

// findCoveredLineEdges traverses the star of edges, maintaing the current location in the result
// area at this node (if any).
// If any L edges are found in the interior of the result, mark them as covered.
func (des *DirectedEdgeStar) findCoveredLineEdges() {
	// Since edges are stored in CCW order around the node,
	// as we move around the ring we move from the right to the left side of the edge

	// Find first DirectedEdge of result area (if any).
	// The interior of the result is on the RHS of the edge,
	// so the start location will be:
	// - INTERIOR if the edge is outgoing
	// - EXTERIOR if the edge is incoming
	startLoc := location.None

	des.edgeMap.WalkInterruptible(func(key, value interface{}) bool {

		nextOut := value.(DirectedEdge)
		nextIn := nextOut.sym
		if !nextOut.isLineEdge() {
			if nextOut.isInResult() {
				startLoc = location.Interior
				return false
			}
			if nextIn.isInResult {
				startLoc = location.Exterior
				return false
			}
		}
		return true
	})
	// no A edges found, so can't determine if L edges are covered or not
	if startLoc == location.None {
		return
	}

	// move around ring, keeping track of the current location
	// (Interior or Exterior) for the result area.
	// If L edges are found, mark them as covered if they are in the interior
	currLoc := startLoc

	des.edgeMap.Walk(func(key, value interface{}) {
		nextOut := value.(DirectedEdge)
		nextIn := nextOut.sym
		if nextOut.isLineEdge() {
			nextOut.edge.isCovered = currLoc == location.Interior
		} else {
			// edge is an Area edge
			if nextOut.isInResult() {
				currLoc = location.Exterior
			}
			if nextIn.isInResult() {
				currLoc = location.Interior
			}
		}
	})
}
func (des *DirectedEdgeStar) computeDepths(de DirectedEdge) {
	edgeIndex := des.findIndex(de)
	startDepth := de.depth[LEFT]
	targetLastDepth := de.depth[RIGHT]
	// compute the depths from this edge up to the end of the edge array
	nextDepth := des.computeDepthsForEdges(edgeIndex+1, len(des.edgeList), startDepth)
	// compute the depths for the initial part of the array
	lastDepth := des.computeDepthsForEdges(0, edgeIndex, nextDepth)
	//Debug.print(lastDepth != targetLastDepth, this);
	//Debug.print(lastDepth != targetLastDepth, "mismatch: " + lastDepth + " / " + targetLastDepth);
	if lastDepth != targetLastDepth {
		return fmt.Errorf("depth mismatch at %v", de.Coordinate())
	}
}

// computeDepths calculates the DirectedEdge depths for a subsequence of the edge array.
func (des *DirectedEdgeStar) computeDepthsForEdges(startIndex, endIndex, startDepth int) int {
	currDepth := startDepth
	for i := startIndex; i < endIndex; i++ {
		nextDe := des.edgeList[i].(DirectedEdge)
		nextDe.setEdgeDepths(RIGHT, currDepth)
		currDepth = nextDe.depth[LEFT]
	}
	return currDepth
}
