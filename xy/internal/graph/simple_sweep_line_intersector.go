package graph

type simpleSweepLineIntersector struct {
	nOverlaps int
	events    []*SweepLineEvent
}

var _ edgeSetIntersector = &simpleSweepLineIntersector{}

func (s *simpleSweepLineIntersector) computeIntersections(edges []*Edge, si SegmentIntersector, testAllSegments bool) {
	if testAllSegments {
		s.addEdgeListToEdgeSet(edges, nil)
	} else {
		s.addEdges(edges)
	}
	s.computeIntersectionsFromEvents(si)
}
func (s *simpleSweepLineIntersector) computeIntersectionsForEdges(edges0, edges1 []*Edge, si SegmentIntersector) {
	s.addEdgeListToEdgeSet(edges0, edges0)
	s.addEdgeListToEdgeSet(edges1, edges1)
	s.computeIntersectionsFromEvents(si)
}

func (s *simpleSweepLineIntersector) addEdges(edges []*Edge) {
	for _, edge := range edges {
		// edge is its own group
		s.addEdgeToEdgeSet(edge, edge)
	}
}

func (s *simpleSweepLineIntersector) addEdgeListToEdgeSet(edges []*Edge, edgeSet interface{}) {
	for _, edge := range edges {
		s.addEdgeToEdgeSet(edge, edgeSet)
	}
}
func (s *simpleSweepLineIntersector) addEdgeToEdgeSet(edge *Edge, edgeSet interface{}) {
	pts := edge.pts

	for i := 0; i < len(pts)-1; i++ {
		ss := NewSweepLineSegment(edge, i)
		insertEvent := NewSweepLineEvent(edgeSet, ss.MinX(), nil, ss)
		s.events = append(s.events, insertEvent)
		s.events = append(s.events, NewSweepLineEvent(edgeSet, ss.MaxX(), insertEvent, ss))
	}
}

func (s *simpleSweepLineIntersector) computeIntersectionsFromEvents(si SegmentIntersector) {
	s.nOverlaps = 0
	prepareEvents(s.events)

	for i, ev := range s.events {
		if ev.eventType == INSERT {
			s.processOverlaps(i, ev.deleteEventIndex, ev, si)
		}
	}
}

func (s *simpleSweepLineIntersector) processOverlaps(start, end int, ev0 *SweepLineEvent, si SegmentIntersector) {
	ss0 := ev0.obj.(*SweepLineSegment)
	/**
	 * Since we might need to test for self-intersections,
	 * include current insert event object in list of event objects to test.
	 * Last index can be skipped, because it must be a Delete event.
	 */
	for i := start; i < end; i++ {
		ev1 := s.events[i]
		if ev1.eventType == INSERT {
			ss1 := ev1.obj.(*SweepLineSegment)
			if ev0.edgeSet == nil || (ev0.edgeSet != ev1.edgeSet) {
				ss0.computeIntersections(ss1, si)
				s.nOverlaps++
			}
		}
	}
}

type SweepLineSegment struct {
	edge    *Edge
	pts     []float64
	ptIndex int
}

func NewSweepLineSegment(edge *Edge, ptIndex int) *SweepLineSegment {
	return &SweepLineSegment{
		edge:    edge,
		ptIndex: ptIndex,
		pts:     edge.pts,
	}
}

func (s *SweepLineSegment) MinX() float64 {
	x1 := s.pts[s.ptIndex*s.edge.layout.Stride()]
	x2 := s.pts[s.ptIndex*s.edge.layout.Stride()+s.edge.layout.Stride()]
	if x1 < x2 {
		return x1
	}
	return x2
}

func (s *SweepLineSegment) MaxX() float64 {
	stride := s.edge.layout.Stride()
	x1 := s.pts[s.ptIndex*stride]
	x2 := s.pts[s.ptIndex*stride+stride]

	if x1 > x2 {
		return x1
	}
	return x2
}

func (s *SweepLineSegment) computeIntersections(ss *SweepLineSegment, si SegmentIntersector) {
	si.addIntersections(s.edge, s.ptIndex, ss.edge, ss.ptIndex)
}
