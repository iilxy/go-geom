package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/location"
)

type edgeRingStrategy interface {
	getNext(dEdge *directedEdge) *directedEdge
	setEdgeRing(dEdge *directedEdge, edgeRing *edgeRing)
}

type edgeRing struct {
	strategy      edgeRingStrategy
	layout        geom.Layout
	startDe       *directedEdge
	maxNodeDegree int
	edges         []*directedEdge
	pts           []float64
	label         *Label
	ring          *geom.LinearRing
	shell         *edgeRing
	holes         []*edgeRing
}

func nNewEdgeRingCommon(layout geom.Layout, start *directedEdge) *edgeRing {
	EdgeRingCommon := &edgeRing{
		layout:        layout,
		maxNodeDegree: -1,
		label:         NewHomogeneousLabel(location.None),
	}

	EdgeRingCommon.computePoints(start)
	EdgeRingCommon.computeRing()
	return EdgeRingCommon
}

func (er *edgeRing) getNext(dEdge *directedEdge) *directedEdge {
	return er.strategy.getNext(dEdge)
}
func (er *edgeRing) setEdgeRing(dEdge *directedEdge, edgeRing *edgeRing) {
	er.strategy.setEdgeRing(dEdge, edgeRing)
}

func (er *edgeRing) isIsolated() bool {
	return er.label.getGeometryCount() == 1
}

func (er *edgeRing) isHole() bool {
	return er.shell != nil
}

func (er *edgeRing) setShell(shell *edgeRing) {
	er.shell = shell
	shell.holes = append(shell.holes, er)
}

func (er *edgeRing) toPolygon() *geom.Polygon {
	ends := make([]int, 1+len(er.holes), 0)
	shellLen := len(er.ring.FlatCoords())
	ends = append(ends, shellLen)
	numOrds := shellLen
	for i := 0; i < len(er.holes); i++ {
		ringLen := len(er.holes[i].ring.FlatCoords())
		numOrds += ringLen
		ends = append(ends, ringLen)
	}

	holeData := make([]float64, numOrds, 0)
	holeData = append(holeData, er.ring.FlatCoords()...)
	for i := 0; i < len(er.holes); i++ {
		holeData = append(holeData, er.holes[i].ring.FlatCoords()...)
	}
	return geom.NewPolygonFlat(er.shell.ring.Layout(), holeData, ends)
}

// Compute a LinearRing from the point list previously collected.
// Test if the ring is a hole (i.e. if it is CCW) and set the hole flag
// accordingly.
func (er *edgeRing) computeRing() {
	if er.ring != nil {
		return
	}
	er.ring = geom.NewLinearRingFlat(er.layout, er.pts)
}

func (er *edgeRing) computePoints(start *directedEdge) error {
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
		er.addPoints(de.edge, de.isForward, isFirstEdge)
		isFirstEdge = false
		er.setEdgeRing(de, er)
		de = er.getNext(de)
		if de == startDe {
			break
		}
	}
	return nil
}

func (er *edgeRing) getMaxNodeDegree() int {
	if er.maxNodeDegree < 0 {
		er.computeMaxNodeDegree()
	}
	return er.maxNodeDegree
}

func (er *edgeRing) computeMaxNodeDegree() {
	er.maxNodeDegree = 0
	de := er.startDe
	for {
		node := de.node
		degree := node.edges.(*directedEdgeStar).getOutgoingDegreeInRing(er)
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

func (er *edgeRing) setInResult() {
	de := er.startDe
	for {
		de.edge.isInResult = true
		de = de.next
		if de == er.startDe {
			return
		}
	}
}

func (er *edgeRing) fullLabelMerge(deLabel *Label) {
	er.mergeLabel(deLabel, 0)
	er.mergeLabel(deLabel, 1)
}

// mergeLabel merges the RHS label from a DirectedEdge into the label for this EdgeRing.
// The DirectedEdge label may be null.  This is acceptable - it results
// from a node which is NOT an intersection node between the Geometries
// (e.g. the end node of a LinearRing).  In this case the DirectedEdge label
// does not contribute any information to the overall labelling, and is simply skipped.
func (er *edgeRing) mergeLabel(deLabel *Label, geomIndex int) {
	loc := deLabel[geomIndex][RIGHT]
	// no information to be had from this label
	if loc == location.None {
		return
	}
	// if there is no current RHS value, set it
	if er.label[geomIndex][ON] == location.None {
		er.label[geomIndex][ON] = loc
		return
	}
}

func (er *edgeRing) addPoints(edge *Edge, isForward, isFirstEdge bool) {
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

func (er *edgeRing) containsPoint(p geom.Coord) bool {
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
