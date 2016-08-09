package xygraph

import (
	"sort"
)

type simpleMCSweepLineIntersector struct {
	nOverlaps int
	events    []*SweepLineEvent
}

var _ edgeSetIntersector = &simpleMCSweepLineIntersector{}

func (s *simpleMCSweepLineIntersector) computeIntersections(edges []*Edge, si SegmentIntersector, testAllSegments bool) {

	if testAllSegments {
		s.addEdgeListToEdgeSet(edges, nil)
	} else {
		s.addEdges(edges)
	}
	s.computeIntersectionsFromEvents(si)
}

func (s *simpleMCSweepLineIntersector) computeIntersectionsForEdges(edges0, edges1 []*Edge, si SegmentIntersector) {
	s.addEdgeListToEdgeSet(edges0, edges0)
	s.addEdgeListToEdgeSet(edges1, edges1)
	s.computeIntersectionsFromEvents(si)
}

func (s *simpleMCSweepLineIntersector) addEdges(edges []*Edge) {
	for _, edge := range edges {
		// edge is its own group
		s.addEdgeToEdgeSet(edge, edge)
	}
}
func (s *simpleMCSweepLineIntersector) addEdgeListToEdgeSet(edges []*Edge, edgeSet interface{}) {
	for _, edge := range edges {
		s.addEdgeToEdgeSet(edge, edgeSet)
	}
}
func (s *simpleMCSweepLineIntersector) addEdgeToEdgeSet(edge *Edge, edgeSet interface{}) {
	mce := edge.mce
	startIndex := mce.startIndex
	for i, _ := range startIndex {
		mc := monotoneChain{mce: mce, chainIndex: i}
		insertEvent := NewSweepLineEvent(edgeSet, mce.minX(i), nil, mc)
		s.events = append(s.events, insertEvent)
		s.events = append(s.events, NewSweepLineEvent(edgeSet, mce.maxX(i), insertEvent, mc))
	}
}

// Because Delete Events have a link to their corresponding Insert event,
// it is possible to compute exactly the range of events which must be
// compared to a given Insert event object.
func prepareEvents(events []*SweepLineEvent) {
	sort.Sort(SortableSweepLineEvents{events})
	for i, ev := range events {
		if ev.eventType == DELETE {
			ev.insertEvent.deleteEventIndex = i
		}
	}
}

func (s *simpleMCSweepLineIntersector) computeIntersectionsFromEvents(si SegmentIntersector) {
	s.nOverlaps = 0
	prepareEvents(s.events)

	for i, ev := range s.events {
		if ev.eventType == INSERT {
			s.processOverlaps(i, ev.deleteEventIndex, ev, si)
		}
	}
}

func (s *simpleMCSweepLineIntersector) processOverlaps(start, end int, ev0 *SweepLineEvent, si SegmentIntersector) {
	mc0 := ev0.obj.(monotoneChain)
	// Since we might need to test for self-intersections,
	// include current insert event object in list of event objects to test.
	// Last index can be skipped, because it must be a Delete event.

	for i := start; i < end; i++ {
		ev1 := s.events[i]
		if ev1.eventType == INSERT {
			mc1 := ev1.obj.(monotoneChain)
			// don't compare edges in same group
			// null group indicates that edges should be compared
			if ev0.edgeSet == nil || ev0.edgeSet != ev1.edgeSet {
				mc0.computeIntersections(mc1, si)
				s.nOverlaps++
			}
		}
	}
}

type SweepLineEventType int

const (
	INSERT SweepLineEventType = iota + 1
	DELETE
)

type SweepLineEvent struct {
	edgeSet          interface{}
	xValue           float64
	eventType        SweepLineEventType
	insertEvent      *SweepLineEvent
	deleteEventIndex int
	obj              interface{}
}

func NewSweepLineEvent(edgeSet interface{}, x float64, insertEvent *SweepLineEvent, obj interface{}) *SweepLineEvent {
	var eventType = DELETE

	if insertEvent == nil {
		eventType = INSERT
	}

	return &SweepLineEvent{
		edgeSet:     edgeSet,
		xValue:      x,
		insertEvent: insertEvent,
		eventType:   eventType,
		obj:         obj,
	}
}

type SortableSweepLineEvents struct {
	events []*SweepLineEvent
}

func (s SortableSweepLineEvents) Len() int {
	return len(s.events)
}
func (s SortableSweepLineEvents) Less(i, j int) bool {
	e1 := s.events[i]
	e2 := s.events[2]
	if e1.xValue < e2.xValue {
		return true
	}
	if e1.xValue > e2.xValue {
		return false
	}
	if e1.eventType < e2.eventType {
		return true
	}
	return false
}

func (s SortableSweepLineEvents) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}
