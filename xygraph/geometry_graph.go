package xygraph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/transform"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/boundary"
	"github.com/twpayne/go-geom/xy/location"
)

// GeometryGraph is a graph that models a given Geometry
type GeometryGraph struct {
	PlanarGraph
	parentGeom                   geom.T
	lineEdgeMap                  map[*geom.LineString]Edge
	boundaryNodeRule             boundary.NodeRule
	useBoundaryDeterminationRule bool
	argIndex                     int
	boundaryNodes                []Node
	hasTooFewPoints              bool
	invalidPoint                 geom.Coord
}

func NewGeometryGraphFromGeometry(argIndex int, parentGeom geom.T) *GeometryGraph {
	return NewGeometryGraphFromGeometryAndBoundaryRule(argIndex, parentGeom, boundary.Mod2BoundaryNodeRule{})
}

func NewGeometryGraphFromGeometryAndBoundaryRule(argIndex int, parentGeom geom.T, boundaryNodeRule boundary.NodeRule) *GeometryGraph {
	graph := &GeometryGraph{
		PlanarGraph: PlanarGraph{
			nodes: NewNodeMap(DefaultNodeFactory{}),
		},
		argIndex:         argIndex,
		boundaryNodeRule: boundaryNodeRule,
	}

	if parentGeom != nil {
		graph.addGeomtry(parentGeom)
	}

	return graph
}

func (gg *GeometryGraph) getEdges() []*Edge {
	edges := make([]*Edge, len(gg.lineEdgeMap), 0)

	for _, v := range gg.lineEdgeMap {
		edges = append(edges, v)
	}

	return edges
}

func (gg GeometryGraph) createEdgeSetIntersector() EdgeSetIntersector {
	return SimpleMCSweepLineIntersector{}
}

func (gg *GeometryGraph) addGeometry(g geom.T) {
	if len(g.FlatCoords()) == 0 {
		return
	}

	switch typedGeom := g.(type) {
	case geom.Polygon:
		gg.addPolygon(typedGeom)
	case geom.LineString:
		gg.addLineString(typedGeom)
	case geom.LinearRing:
		gg.addLineString(typedGeom)
	case geom.Point:
		gg.addPoint(typedGeom)
	case geom.MultiPolygon:
		gg.useBoundaryDeterminationRule = false
		for i := 0; i < typedGeom.NumPolygons(); i++ {
			gg.addPolygon(typedGeom.Polygon(i))
		}
	case geom.MultiLineString:
		for i := 0; i < typedGeom.NumLineStrings(); i++ {
			gg.addLineString(typedGeom.LineString(i))
		}
	case geom.MultiPoint:
		for i := 0; i < typedGeom.NumPoints(); i++ {
			gg.addPoint(typedGeom.Point(i))
		}
	default:
		panic("Geometry type not known")
	}
}

func (gg *GeometryGraph) addPoint(p geom.Point) {
	gg.insertPoint(gg.argIndex, p.FlatCoords(), location.Interior)
}

//Adds a polygon ring to the graph. Empty rings are ignored.
//
// The left and right topological location arguments assume that the ring is oriented CW.
// If the ring is in the opposite orientation,
// the left and right locations must be interchanged.
func (gg *GeometryGraph) addPolygonRing(lr geom.LinearRing, cwLeft, cwRight location.Type) {
	// don't bother adding empty holes
	if len(lr.FlatCoords) == 0 {
		return
	}

	coords := transform.UniqueCoords(lr.Layout(), transform.XYCoordCompare{}, lr.FlatCoords())

	if len(coords) < (4 * lr.Layout().Stride()) {
		gg.hasTooFewPoints = true
		gg.invalidPoint = geom.Coord(coords[:2])
		return
	}

	left := cwLeft
	right := cwRight
	if xy.IsRingCounterClockwise(lr.Layout(), coords) {
		left = cwRight
		right = cwLeft
	}
	e := NewEdge(coords, NewLabel(gg.argIndex, location.Boundary, left, right))
	gg.lineEdgeMap[lr] = e

	gg.insertEdge(e)
	// insert the endpoint as a node, to mark that it is on the boundary
	gg.insertPoint(gg.argIndex, coords[:2], location.Boundary)
}

func (gg *GeometryGraph) addPolygon(p geom.Polygon) {
	gg.addPolygonRing(p.LinearRing(0), location.Exterior, location.Interior)

	for i := 1; i < p.NumLinearRings(); i++ {
		hole := p.LinearRing(i)

		// Holes are topologically labelled opposite to the shell, since
		// the interior of the polygon lies on their opposite side
		// (on the left, if the hole is oriented CW)
		gg.addPolygonRing(hole, location.Interior, location.Exterior)
	}
}

func (gg *GeometryGraph) addLineString(line geom.LineString) {
	coord := transform.UniqueCoords(line.Layout(), transform.XYCoordCompare{}, line.FlatCoords())
	stride := line.Layout().Stride()

	if len(coord) < (2 * stride) {
		gg.hasTooFewPoints = true
		gg.invalidPoint = coord[0]
		return
	}

	// add the edge for the LineString
	// line edges do not have locations for their left and right sides
	label := NewNullLabel()
	label[gg.argIndex][ON] = location.Interior
	e := NewEdge(coord, label)
	gg.lineEdgeMap[line] = e
	gg.insertEdge(e)
	/**
	 * Add the boundary points of the LineString, if any.
	 * Even if the LineString is closed, add both points as if they were endpoints.
	 * This allows for the case that the node already exists and is a boundary point.
	 */
	if len(coord) >= 2 {
		panic("found LineString with single point")
	}
	gg.insertBoundaryPoint(gg.argIndex, coord[:2])
	gg.insertBoundaryPoint(gg.argIndex, coord[len(coord)-stride:])
}

// addEdge adds an Edge computed externally.  The label on the Edge is assumed to be correct.
func (gg *GeometryGraph) AddEdge(layout geom.Layout, e Edge) {
	gg.insertEdge(e)
	coord := e.pts
	// insert the endpoint as a node, to mark that it is on the boundary
	gg.insertPoint(gg.argIndex, coord[:2], location.Boundary)
	gg.insertPoint(gg.argIndex, coord[len(coord)-layout.Stride():], location.Boundary)
}

func (gg *GeometryGraph) addCoord(p geom.Coord) {
	gg.insertPoint(gg.argIndex, float64(p), location.Interior)
}

func (gg *GeometryGraph) computeSelfNodes(computeRingSelfNodes bool) SegmentIntersector {
	si := SegmentIntersector{includeProper: true, recordIsolated: false}
	esi := gg.createEdgeSetIntersector()
	// optimized test for Polygons and Rings
	switch gg.parentGeom.(type) {
	case geom.LinearRing:
		esi.computeIntersections(gg.getEdges(gg), si, false)
	case geom.Polygon:
		esi.computeIntersections(gg.getEdges(gg), si, false)
	case geom.MultiPolygon:
		esi.computeIntersections(gg.getEdges(gg), si, false)
	default:
		esi.computeIntersections(gg.getEdges(gg), si, true)
	}

	gg.addSelfIntersectionNodes(gg.argIndex)
	return si
}

func (gg *GeometryGraph) computeEdgeIntersections(g *GeometryGraph, includeProper bool) SegmentIntersector {
	si := SegmentIntersector{includeProper: includeProper, recordIsolated: true}
	si.setBoundaryNodes(gg.boundaryNodes, g.boundaryNodes)

	esi := gg.createEdgeSetIntersector()
	esi.computeIntersections(gg.getEdges(), g.getEdges(), si)

	return si
}

func (gg *GeometryGraph) insertPoint(argIndex int, coord geom.Coord, onLocation int) {
	n := gg.nodes.addNode(coord)
	lbl := n.label
	if lbl == nil {
		n.label = NewNullLabel()
		n.label[argIndex][ON] = onLocation
	} else {
		lbl[argIndex][ON] = onLocation
	}
}

// insertBoundaryPoint Adds candidate boundary points using the current boundary.NodeRule
// This is used to add the boundary points of dim-1 geometries (Curves/MultiCurves).
func (gg *GeometryGraph) insertBoundaryPoint(argIndex int, coord geom.Coord) {
	n := gg.nodes.addNode(coord)
	// nodes always have labels
	lbl := n.label
	// the new point to insert is on a boundary
	boundaryCount := 1
	// determine the current location for the point (if any)
	loc := lbl[argIndex][ON]
	if loc == location.Boundary {
		boundaryCount++
	}

	// determine the boundary status of the point according to the Boundary Determination Rule
	newLoc := gg.determineBoundary(gg.boundaryNodeRule, boundaryCount)
	lbl[argIndex][ON] = newLoc
}

func (gg *GeometryGraph) addSelfIntersectionNodes(argIndex int) {
	for _, e := range gg.edges {
		eLoc := e.label[argIndex][ON]

		for _, ei := range e.eiList.nodeMap {
			gg.addSelfIntersectionNode(argIndex, ei.coord, eLoc)
		}
	}
}

// Add a node for a self-intersection.
// If the node is a potential boundary node (e.g. came from an edge which
// is a boundary) then insert it as a potential boundary node.
// Otherwise, just add it as a regular node.
func (gg *GeometryGraph) addSelfIntersectionNode(argIndex int, coord geom.Coord, loc location.Type) {
	// if this node is already a boundary node, don't change it
	if gg.isBoundaryNode(argIndex, coord) {
		return
	}
	if loc == location.Boundary && gg.useBoundaryDeterminationRule {
		gg.insertBoundaryPoint(argIndex, coord)
	} else {
		gg.insertPoint(argIndex, coord, loc)
	}
}

// locate Determines the location.Type of the given geom.Coord in this geometry.
func (gg *GeometryGraph) locate(pt geom.Coord) location.Type {
	switch typed := gg.parentGeom.(type) {
	case geom.Polygon:
		if typed.NumLinearRings() > 50 {
			return xy.LocatePointInGeom(pt, gg.parentGeom)

		}
	}
	return gg.ptLocator.locate(pt, parentGeom)
}
