package xygraph

import "github.com/twpayne/go-geom"

type MonotoneChainEdge struct {
	e          Edge
	pts        []geom.Coord
	startIndex []int
	env1, env2 geom.Bounds
}

func newMonotoneChainEdge(e Edge) *MonotoneChainEdge {
	return &MonotoneChainEdge{
		e:          e,
		pts:        e.pts,
		startIndex: getChainStartIndices(e.pts),
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
