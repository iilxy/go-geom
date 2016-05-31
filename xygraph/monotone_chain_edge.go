package xygraph

import "github.com/twpayne/go-geom"

type MonotoneChainEdge struct {
	e          *Edge
	pts        []geom.Coord
	startIndex []int
	env1, env2 geom.Bounds
}

func newMonotoneChainEdge(e *Edge) *MonotoneChainEdge {
	return &MonotoneChainEdge{
		e:          e,
		pts:        e.pts,
		startIndex: getChainStartIndices(e.pts),
	}
}

func (mce *MonotoneChainEdge) MinX(chainIndex int) float64 {
	x1 := mce.pts[mce.startIndex[chainIndex]][0]
	x2 := mce.pts[mce.startIndex[chainIndex+1]][0]
	if x1 < x2 {
		return x1
	}
	return x2
}

func (mce *MonotoneChainEdge) MaxX(chainIndex int) float64 {
	x1 := mce.pts[mce.startIndex[chainIndex]][0]
	x2 := mce.pts[mce.startIndex[chainIndex+1]][0]
	if x1 > x2 {
		return x1
	}
	return x2
}
func (mce *MonotoneChainEdge) computeIntersectsForChain(chainIndex0 int, otherMCE MonotoneChainEdge, chainIndex1 int, si SegmentIntersector) {
	otherMCE.computeIntersectsForChainBounded(
		mce.startIndex[chainIndex0], mce.startIndex[chainIndex0+1],
		otherMCE,
		otherMCE.startIndex[chainIndex1], otherMCE.startIndex[chainIndex1+1],
		si)
}
func (mce *MonotoneChainEdge) computeIntersectsForChainBounded(start0, end0 int, otherMCE MonotoneChainEdge, start1, end1 int, ei SegmentIntersector) {
	p00 := mce.pts[start0]
	p01 := mce.pts[end0]
	p10 := otherMCE.pts[start1]
	p11 := otherMCE.pts[end1]

	// terminating condition for the recursion
	if end0-start0 == 1 && end1-start1 == 1 {
		ei.addIntersections(mce.e, start0, otherMCE.e, start1)
		return
	}
	// nothing to do if the envelopes of these chains don't overlap
	mce.env1.SetCoords(p00, p01)
	mce.env2.SetCoords(p10, p11)
	if !mce.env1.Overlaps(mce.env2) {
		return
	}

	// the chains overlap, so split each in half and iterate  (binary search)
	mid0 := (start0 + end0) / 2
	mid1 := (start1 + end1) / 2

	// Assert: mid != start or end (since we checked above for end - start <= 1)
	// check terminating conditions before recursing
	if start0 < mid0 {
		if start1 < mid1 {
			mce.computeIntersectsForChainBounded(start0, mid0, otherMCE, start1, mid1, ei)
		}
		if mid1 < end1 {
			mce.computeIntersectsForChainBounded(start0, mid0, otherMCE, mid1, end1, ei)
		}
	}
	if mid0 < end0 {
		if start1 < mid1 {
			mce.computeIntersectsForChainBounded(mid0, end0, otherMCE, start1, mid1, ei)
		}
		if mid1 < end1 {
			mce.computeIntersectsForChainBounded(mid0, end0, otherMCE, mid1, end1, ei)
		}
	}
}

func getChainStartIndices(pts []geom.Coord) []int {
	// find the startpoint (and endpoints) of all monotone chains in this edge
	start := 0
	startIndexList := []int{start}
	for {
		last := findChainEnd(pts, start)
		startIndexList = append(startIndexList, last)
		start = last
		if !(start < len(pts)-1) {
			break
		}
	}

	return startIndexList
}

func findChainEnd(pts []geom.Coord, start int) int {
	// determine quadrant for chain
	chainQuad := coordsQuadrant(pts[start], pts[start+1])
	last := start + 1
	for last < len(pts) {
		// compute quadrant for next possible segment in chain
		quad := coordsQuadrant(pts[last-1], pts[last])
		if quad != chainQuad {
			break
		}
		last++
	}
	return last - 1
}
