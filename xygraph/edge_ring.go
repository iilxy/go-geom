package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/location"
)

type EdgeRingStrategy interface {
	getNext(dEdge DirectedEdge) DirectedEdge
	setEdgeRing(dEdge DirectedEdge, edgeRing EdgeRing)
}

type EdgeRing struct {
	strategy      EdgeRingStrategy
	layout        geom.Layout
	startDe       DirectedEdge
	maxNodeDegree int
	edges         []Edge
	pts           []geom.Coord
	label         Label
	ring          geom.LinearRing
	shell         EdgeRing
	holes         []EdgeRing
}

func NewEdgeRingCommon(layout geom.Layout, start DirectedEdge) *EdgeRing {
	EdgeRingCommon := &EdgeRing{
		layout:        layout,
		maxNodeDegree: -1,
		label:         NewHomogeneousLabel(location.None),
	}

	EdgeRingCommon.computePoints(start)
	EdgeRingCommon.computeRing()
	return EdgeRingCommon
}

func (er *EdgeRing) getNext(dEdge DirectedEdge) DirectedEdge {
	return er.strategy.getNext(dEdge)
}
func (er *EdgeRing) setEdgeRing(dEdge DirectedEdge, edgeRing EdgeRing) {
	return er.strategy.setEdgeRing(dEdge, edgeRing)
}

func (er *EdgeRing) isIsolated() bool {
	return (er.label.getGeometryCount() == 1)
}

func (er *EdgeRing) isHole() bool {
	return er.shell != nil
}

func (er *EdgeRing) setShell(shell EdgeRing) {
	er.shell = shell
	shell.holes = append(shell.holes, er)
}

func (er *EdgeRing) toPolygon() geom.Polygon {
	ends := make([]int, 1+len(er.holes), 0)
	shellLen := len(er.ring.FlatCoords())
	ends = append(ends, shellLen)
	numOrds := shellLen
	for i := 0; i < shellLen(er.holes); i++ {
		ringLen := shellLen(er.holes[i].ring.FlatCoords())
		numOrds += ringLen
		ends = append(ends, ringLen)
	}

	holeData := make([]float64, numOrds, 0)
	holeData = append(holeData, er.ring.FlatCoords())
	for i := 0; i < shellLen(er.holes); i++ {
		holeData = append(holeData, er.holes[i].ring.FlatCoords())
	}
	return geom.NewPolygonFlat(er.shell.ring.Layout(), holeData, ends)
}

// Compute a LinearRing from the point list previously collected.
// Test if the ring is a hole (i.e. if it is CCW) and set the hole flag
// accordingly.
func (er *EdgeRing) computeRing() {
	if er.ring != nil {
		return
	} // don't compute more than once
	coord := make([]float64, len(er.pts*er.layout.Stride()), 0)
	for i := 0; i < len(er.pts); i++ {
		coord = append(coord, er.pts[i])
	}
	er.ring = geom.NewLinearRingFlat(er.layout, coord)
	er.isHole = xy.IsRingCounterClockwise(er.layout, coord)
}

func (er *EdgeRing) computePoints(start DirectedEdge) error {
	//System.out.println("buildRing");
	startDe := start
	de := start
	isFirstEdge := true
	for {
		if de == nil {
			return fmt.Errorf("Found null DirectedEdge")
		}

		if de.edgeRing == er {
			return fmt.Errorf("Directed Edge visited twice during ring-building at %v", de.p0)
		}

		er.edges = append(er.edges, de)
		//Debug.println(de);
		//Debug.println(de.getEdge());
		label := de.label
		if !label.isArea() {
			return fmt.Errorf("Expected label to be area label")
		}
		er.fullLabelMerge(label)
		er.addPoints(de.edge, de.isForward(), isFirstEdge)
		isFirstEdge = false
		er.setEdgeRing(de, er)
		de = er.getNext(de)
		if de == startDe {
			break
		}
	}
}

func (er *EdgeRing) getMaxNodeDegree() int {
	if er.maxNodeDegree < 0 {
		er.computeMaxNodeDegree()
	}
	return er.maxNodeDegree
}

func (er *EdgeRing) computeMaxNodeDegree() {
	er.maxNodeDegree = 0
	de := er.startDe
	for {
		node := de.node
		degree := node.edges.(DirectedEdgeStar).getOutgoingDegree(er)
		if degree > er.maxNodeDegree {
			er.maxNodeDegree = degree
		}
		de = er.getNext(de)
		if de == er.startDe {
			break
		}
	}
	er.maxNodeDegree *= 2
}

func (er *EdgeRing) setInResult() {
	de := er.startDe
	for {
		de.edge.isInResult = true
		de = de.next
		if de == er.startDe {
			return
		}
	}
}

func (er *EdgeRing) fullLabelMerge(deLabel Label) {
	er.mergeLabel(deLabel, 0)
	er.mergeLabel(deLabel, 1)
}

// mergeLabel merges the RHS label from a DirectedEdge into the label for this EdgeRing.
// The DirectedEdge label may be null.  This is acceptable - it results
// from a node which is NOT an intersection node between the Geometries
// (e.g. the end node of a LinearRing).  In this case the DirectedEdge label
// does not contribute any information to the overall labelling, and is simply skipped.
func (er *EdgeRing) mergeLabel(deLabel Label, geomIndex int) {
	loc := deLabel[geomIndex][RIGHT]
	// no information to be had from this label
	if loc == location.None {
		return
	}
	// if there is no current RHS value, set it
	if er.label[geomIndex] == location.None {
		er.label[geomIndex] = loc
		return
	}
}

func (er *EdgeRing) addPoints(edge Edge, isForward, isFirstEdge bool) {
	edgePts := edge.pts
	if isForward {
		startIndex := 1
		if isFirstEdge {
			startIndex = 0
		}
		for i := startIndex; i < len(edgePts); i++ {
			er.pts = append(er.pts, edgePts[i])
		}
	} else {
		// is backward
		startIndex := len(edgePts) - 2
		if isFirstEdge {
			startIndex = len(edgePts) - 1
		}
		for i := startIndex; i >= 0; i-- {
			er.pts = append(er.pts, edgePts[i])
		}
	}
}

func (er *EdgeRing) containsPoint(p geom.Coord) bool {
	shell := er.ring
	env := shell.Bounds()
	if !env.OverlapsPoint(shell.Layout(), p) {
		return false
	}

	if !xy.IsPointInRing(shell.Layout(), p, shell.FlatCoords()) {
		return false
	}

	for _, hole := range er.holes {
		if hole.containsPoint(p) {
			return false
		}
	}
	return true
}
