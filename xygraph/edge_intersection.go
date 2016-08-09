package xygraph

import (
	"bufio"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

type edgeIntersectionId [2]float64
type edgeIntersection struct {
	id           edgeIntersectionId
	coord        geom.Coord
	segmentIndex int
	dist         float64
}

func newEdgeIntersection(coord geom.Coord, segmentIndex int, dist float64) *edgeIntersection {
	id := edgeIntersectionId{coord[0], coord[1]}

	return &edgeIntersection{
		id: id, coord: coord, segmentIndex: segmentIndex, dist: dist,
	}
}

// compare compares this intersection to the location indicated by the parameters
// Returns -1 this EdgeIntersection is located before the argument location
// Returns 0 this EdgeIntersection is at the argument location
// Returns 1 this EdgeIntersection is located after the argument location
func (ei *edgeIntersection) compare(segmentIndex int, dist float64) int {
	if ei.segmentIndex < segmentIndex {
		return -1
	}
	if ei.segmentIndex > segmentIndex {
		return 1
	}
	if ei.dist < dist {
		return -1
	}
	if ei.dist > dist {
		return 1
	}
	return 0
}
func (ei *edgeIntersection) print(out bufio.Writer) {
	out.WriteString(fmt.Sprintf("%v seg # = %v dist = %v\n", ei.coord, ei.segmentIndex, ei.dist))
}

type edgeIntersectionList struct {
	edge    *Edge
	// key is id of an EdgeIntersection
	nodeMap map[edgeIntersectionId]*edgeIntersection
}

func (ei *edgeIntersectionList) add(intPt geom.Coord, segmentIndex int, dist float64) *edgeIntersection {
	eiNew := newEdgeIntersection(intPt, segmentIndex, dist)
	if eiOld, found := ei.nodeMap[eiNew.id]; found {
		return eiOld
	} else {
		ei.nodeMap[eiNew.id] = eiNew
		return eiNew
	}
}

func (ei *edgeIntersectionList) isIntersection(pt geom.Coord) bool {
	for _, ei := range ei.nodeMap {
		if xy.Equal(ei.coord, 0, pt, 0) {
			return true
		}
	}
	return false
}

func (ei *edgeIntersectionList) addEndpoints() {
	stride := ei.edge.layout.Stride()
	maxSegIndex := len(ei.edge.pts) - stride
	ei.add(geom.Coord(ei.edge.pts[:stride]), 0, 0.0)
	ei.add(ei.edge.pts[maxSegIndex:], maxSegIndex / stride, 0.0)
}

// addSplitEdgesTo creates new edges for all the edges that the intersections in this list split the parent edge into.
//
// Adds the edges to the input list (this is so a single list can be used to accumulate all split
// edges for a Geometry).
func (ei *edgeIntersectionList) addSplitEdgesTo(edgeList []*Edge) []*Edge {
	// ensure that the list has entries for the first and last point of the edge
	ei.addEndpoints()
	var eiPrev *edgeIntersection
	eiPrevInit := false
	for _, eiCurr := range ei.nodeMap {
		if !eiPrevInit {
			eiPrevInit = true
		} else {
			newEdge := ei.createSplitEdge(eiPrev, eiCurr)
			edgeList = append(edgeList, newEdge)
		}
		eiPrev = eiCurr
	}

	return edgeList
}

// createSplitEdge create a new "split edge" with the section of points between
// (and including) the two intersections.
// The label for the new edge is the same as the label for the parent edge.
func (ei *edgeIntersectionList) createSplitEdge(ei0, ei1 *edgeIntersection) *Edge {
	stride := ei.edge.layout.Stride()
	e0StartIndexInArray := ei0.segmentIndex * stride
	e1StartIndexInArray := ei1.segmentIndex * stride
	npts := (ei1.segmentIndex - ei0.segmentIndex + 2) * stride

	// if the last intersection point is not equal to the its segment start pt,
	// add it to the points list as well.
	// (This check is needed because the distance metric is not totally reliable!)
	// The check for point equality is 2D only - Z values are ignored
	useIntPt1 := ei1.dist > 0.0 || !xy.Equal(ei1.coord, 0, ei.edge.pts, e1StartIndexInArray)

	if !useIntPt1 {
		npts-=stride
	}

	pts := make([]float64, npts, 0)
	pts = append(pts, ei0.coord...)
	pts = append(pts, ei.edge.pts[e0StartIndexInArray + stride : e1StartIndexInArray]...)

	if useIntPt1 {
		pts = append(pts, ei1.coord...)
	}

	return NewEdge(ei.edge.layout, pts, NewLabelFromTemplate(ei.edge.label))
}

func (ei *edgeIntersectionList) print(out bufio.Writer) {
	out.WriteString("Intersections:")
	for _, eiCur := range ei.nodeMap {
		eiCur.print(out)
	}
}
