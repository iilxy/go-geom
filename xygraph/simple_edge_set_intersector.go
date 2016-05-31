package xygraph

type SimpleEdgeSetIntersector struct {
}

var _ EdgeSetIntersector = &SimpleEdgeSetIntersector{}

func (s *SimpleEdgeSetIntersector) computeIntersections(edges []*Edge, si SegmentIntersector, testAllSegments bool) {
	for _, edge0 := range edges {
		for _, edge1 := range edges {
			if testAllSegments || edge0 != edge1 {
				s.computeIntersects(edge0, edge1, si)
			}
		}
	}
}
func (s *SimpleEdgeSetIntersector) computeIntersectionsForEdges(edges0, edges1 []*Edge, si SegmentIntersector) {
	for _, edge0 := range edges0 {
		for _, edge1 := range edges1 {
			s.computeIntersects(edge0, edge1, si)
		}
	}
}

func (s *SimpleEdgeSetIntersector) computeIntersects(e0, e1 *Edge, si SegmentIntersector) {
	pts0 := e0.pts
	pts1 := e1.pts

	for i0 := 0; i0 < len(pts0)-1; i0++ {
		for i1 := 0; i1 < len(pts1)-1; i1++ {
			si.addIntersections(e0, i0, e1, i1)
		}
	}
}
