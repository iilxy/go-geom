package graph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/transform"
	"github.com/twpayne/go-geom/xy"
)

// Edge represents an edge in a graph
type Edge struct {
	GraphableImpl
	Intersections edgeIntersections
	IsIsolated    bool
}

var _ Graphable = Edge{}

type EdgeIntersection struct {
	// SegmentIndex is the index of the Edge LineString on which the intersection lies
	SegmentIndex int
	// Coord is the location of the intersection
	Coord geom.Coord
	// Dist is the distance along the line segment (indexed by SegmentIndex) of the intersection (Coord).
	Dist float64
	// The Edge the intersection lies on
	edge *Edge
}

// IsEndPoint checks if the intersection is one of the endpoints of the edge (end thus intersects by definition with
// the other edges coming into that node)
// maxSegmentIndex is the maximum number of segments in the edge (number of LineStrings in the edge)
func (ei EdgeIntersection) IsEndPoint(maxSegmentIndex int) {
	return (ei.SegmentIndex == 0 && ei.Dist == 0.0) || ei.SegmentIndex == maxSegmentIndex
}

type edgeIntersectionCompare struct{}

func (c edgeIntersectionCompare) IsEquals(o1, o2 interface{}) bool {
	i1, i2 := o1.(EdgeIntersection), o2.(EdgeIntersection)
	return i1.Dist == i2.Dist &&
		xy.Equal(i1.Coord, 0, i2.Coord, 0) &&
		&i1.edge == &i2.edge &&
		i1.SegmentIndex == i2.SegmentIndex
}
func (c edgeIntersectionCompare) IsLess(o1, o2 interface{}) bool {
	i1, i2 := o1.(EdgeIntersection), o2.(EdgeIntersection)
	if i1.SegmentIndex < i2.SegmentIndex {
		return true
	} else if i1.SegmentIndex == i2.SegmentIndex {
		return i1.Dist < i2.Dist
	}
	return false
}

// edgeIntersections contains all the intersections this edge has with other edges
type edgeIntersections struct {
	// Edge is the edge that the intersections belong to
	edge *Edge
	// intersections an ordered set of all the intersections.  The first intersections are the ones closest
	// to the edge origin end-point
	intersections *transform.TreeSet
}

func (i *edgeIntersections) checkIntersections() {
	if i.intersections == nil {
		i.intersections = transform.NewTreeSet(edgeIntersectionCompare{})
	}
}

// Add an intersection
func (i *edgeIntersections) Add(intPoint geom.Coord, segmentIndex int, dist float64) {
	i.checkIntersections()
	ei := &EdgeIntersection{
		SegmentIndex: segmentIndex,
		Coord:        intPoint,
		Dist:         dist,
		edge:         i.edge,
	}

	actual, has := i.intersections.Find(ei)
	if has {
		return actual.(*EdgeIntersection)
	}
	return ei
}

// Walk passes each element in the map to the visitor.  The order of visiting is from the element with the smallest key
// to the element with the largest key
func (i *edgeIntersections) Walk(visitor func(intersection EdgeIntersection)) {
	i.checkIntersections()
	i.checkIntersections()
	i.intersections.Walk(visitor())
}

// WalkInterruptible passes each element in the map to the visitor until false is returned from visitor.
// The order of visiting is from the element with the smallest key to the element with the largest key
func (i *edgeIntersections) WalkInterruptible(visitor func(intersection EdgeIntersection) bool) {
	i.checkIntersections()
	i.intersections.WalkInterruptible(visitor)
}

// IsIntersection checks if the coordinate is one of the intersections on the
func (i *edgeIntersections) IsIntersection(coord geom.Coord) (intersects bool) {
	i.WalkInterruptible(func(intersection EdgeIntersection) bool {
		if xy.Equal(intersection.Coord, 0, coord, 0) {
			intersects = true
			return false
		}
		return true
	})

	return intersects
}

func (i *edgeIntersections) addEndpoints() {
	line := i.edge.RelatedObjects[0].Geom
	i.Add(geom.Coord(line.FlatCoords()[:1]), 0, 0.0)
	maxSegIndex := len(line.Ends()) - 1
	i.Add(geom.Coord(line.FlatCoords()[maxSegIndex:]), maxSegIndex, 0.0)
}

func (i *edgeIntersections) addSplitEdges(edgeList []Edge) []Edge {
	// ensure that the list has entries for the first and last point of the edge
	i.addEndpoints()
	var eiPrev *Edge
	i.Walk(func(ei EdgeIntersection) {
		newEdge := i.createSplitEdge(eiPrev, ei)
		edgeList = append(edgeList, newEdge)

		eiPrev = ei
	})

	return edgeList
}

func (i *edgeIntersections) createSplitEdge(ei0, ei1 EdgeIntersection) Edge {
	npts := ei1.SegmentIndex - ei0.SegmentIndex + 2
	relatedGeom := ei0.edge.RelatedObjects[0].Geom
	stride := relatedGeom.Stride()
	lastSegStartPt := relatedGeom.FlatCoords()[ei1.SegmentIndex : ei1.SegmentIndex*stride+stride]
	// if the last intersection point is not equal to the its segment start pt,
	// add it to the points list as well.
	// (This check is needed because the distance metric is not totally reliable!)
	// The check for point equality is 2D only - Z values are ignored
	useIntPt1 := ei1.Dist > 0.0 || !xy.Equal([]float64(ei1.Coord), 0, lastSegStartPt, 0)
	if !useIntPt1 {
		npts--
	}

	pts := make([]float64, npts, 0)
	pts = append(pts, []float64(ei0.Coord)...)

	for i := (ei0.SegmentIndex + 1) * stride; i <= ei1.SegmentIndex; i += stride {
		pts = append([]float64{}, []float64(relatedGeom.FlatCoords()[i:i+stride])...)
	}

	if useIntPt1 {
		pts = append(pts, []float64(ei1.Coord)...)
	}

	return Edge{
		GraphableImpl{
			RelatedObjects: []RelatedObject{
				{
					Geom:  geom.NewLineStringFlat(relatedGeom.Layout(), pts),
					Left:  ei0.edge.RelatedObjects[0].Left,
					On:    ei0.edge.RelatedObjects[0].On,
					Right: ei0.edge.RelatedObjects[0].Right,
				},
			},
		}}
}
