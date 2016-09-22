package graph

import "github.com/twpayne/go-geom"

type monotoneChainEdge struct {
	e          *Edge
	pts        []float64
	startIndex []int
	env1, env2 *geom.Bounds
}

func newMonotoneChainEdge(e *Edge) *monotoneChainEdge {
	return &monotoneChainEdge{
		e:          e,
		pts:        e.pts,
		startIndex: getChainStartIndices(e.layout, e.pts),
	}
}

func (mce *monotoneChainEdge) minX(chainIndex int) float64 {
	x1 := mce.pts[mce.startIndex[chainIndex]]
	x2 := mce.pts[mce.startIndex[chainIndex+mce.e.layout.Stride()]]
	if x1 < x2 {
		return x1
	}
	return x2
}

func (mce *monotoneChainEdge) maxX(chainIndex int) float64 {
	x1 := mce.pts[mce.startIndex[chainIndex]]
	x2 := mce.pts[mce.startIndex[chainIndex+mce.e.layout.Stride()]]
	if x1 > x2 {
		return x1
	}
	return x2
}
func (mce *monotoneChainEdge) computeIntersectsForChain(chainIndex0 int, otherMCE *monotoneChainEdge, chainIndex1 int, si SegmentIntersector) {
	otherMCE.computeIntersectsForChainBounded(
		mce.startIndex[chainIndex0], mce.startIndex[chainIndex0+mce.e.layout.Stride()],
		otherMCE,
		otherMCE.startIndex[chainIndex1], otherMCE.startIndex[chainIndex1+mce.e.layout.Stride()],
		si)
}
func (mce *monotoneChainEdge) computeIntersectsForChainBounded(start0, end0 int, otherMCE *monotoneChainEdge, start1, end1 int, ei SegmentIntersector) {
	stride0 := mce.e.layout.Stride()
	p00 := geom.Coord(mce.pts[start0 : start0+stride0])
	p01 := geom.Coord(mce.pts[end0 : end0+stride0])

	stride1 := otherMCE.e.layout.Stride()
	p10 := geom.Coord(otherMCE.pts[start1 : start1+stride1])
	p11 := geom.Coord(otherMCE.pts[end1 : end1+stride1])

	// terminating condition for the recursion
	if end0-start0 == mce.e.layout.Stride() && end1-start1 == mce.e.layout.Stride() {
		ei.addIntersections(mce.e, start0, otherMCE.e, start1)
		return
	}
	// nothing to do if the envelopes of these chains don't overlap
	mce.env1.SetCoords(p00, p01)
	mce.env2.SetCoords(p10, p11)
	if !mce.env1.Overlaps(geom.XY, mce.env2) {
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

func getChainStartIndices(layout geom.Layout, pts []float64) []int {
	stride := layout.Stride()
	// find the startpoint (and endpoints) of all monotone chains in this edge
	start := 0
	startIndexList := []int{start}
	for {
		last := findChainEnd(stride, pts, start)
		startIndexList = append(startIndexList, last)
		start = last
		if !(start < len(pts)-stride) {
			break
		}
	}

	return startIndexList
}

func findChainEnd(stride int, pts []float64, start int) int {
	// determine quadrant for chain
	chainQuad := coordsQuadrant(geom.Coord(pts[start:start+stride]), geom.Coord(pts[start+stride:start+stride+stride]))
	last := start + stride
	for last < len(pts) {
		// compute quadrant for next possible segment in chain
		quad := coordsQuadrant(geom.Coord(pts[last-stride:last]), geom.Coord(pts[last:last+stride]))
		if quad != chainQuad {
			break
		}
		last += stride
	}
	return last - stride
}
