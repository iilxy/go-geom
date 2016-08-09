package xygraph

import (
	"bytes"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/transform"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/boundary"
	"github.com/twpayne/go-geom/xy/location"
)

type edgeEndStar interface {
	insert(e edgeEnd)
	insertEdgeEnd(e edgeEnd, obj interface{})
	coordinate() geom.Coord
	degree() int
	nextCW(e edgeEnd) edgeEnd
	computeLabelling(geomGraph []GeometryGraph)
	isAreaLabelsConsistent(geomGraph *GeometryGraph) (bool, error)
	findIndex(eSearch edgeEnd) int
	iterate(func(edge edgeEnd) bool)
}
type edgeEndStarCommon struct {
	edgeMap          *transform.TreeMap
	ptInAreaLocation [2]location.Type
}

var _ edgeEndStar = &edgeEndStarCommon{}

func newEdgeEndStarCommon() *edgeEndStarCommon {
	return &edgeEndStarCommon{
		ptInAreaLocation: [2]location.Type{location.None, location.None},
		edgeMap:          transform.NewTreeMap(EdgeEndCompare{}),
	}
}
func (ees *edgeEndStarCommon) iterate(iterFunc func(edge edgeEnd) bool) {
	ees.edgeMap.WalkInterruptible(func(key, value interface{}) bool {
		return iterFunc(value.(edgeEnd))
	})
}
func (ees *edgeEndStarCommon) insert(e edgeEnd) {
	ees.insertEdgeEnd(e, e.(*directedEdge))
}
func (ees *edgeEndStarCommon) insertEdgeEnd(e edgeEnd, obj interface{}) {
	ees.edgeMap.Insert(e, obj)
}

func (ees *edgeEndStarCommon) coordinate() geom.Coord {
	var e edgeEnd = nil

	ees.edgeMap.WalkInterruptible(func(key, value interface{}) bool {
		e = value.(edgeEnd)
		return false
	})

	if e == nil {
		return nil
	}
	return e.getCoord()
}

func (ees *edgeEndStarCommon) degree() int {
	return ees.edgeMap.Size()
}
func (ees *edgeEndStarCommon) nextCW(ee edgeEnd) edgeEnd {
	i := ees.findIndex(ee)
	iNextCW := i - 1
	if i == 0 {
		iNextCW = ees.edgeMap.Size() - 1
	}

	return ees.getEdgeEnd(iNextCW)
}

func (ees *edgeEndStarCommon) getEdgeEnd(index int) edgeEnd {
	var ee edgeEnd = nil
	i := 0
	ees.edgeMap.WalkInterruptible(func(key, value interface{}) bool {

		if index == i {
			ee = value.(edgeEnd)
			return false
		}
		i++
		return true
	})

	return ee
}
func (ees *edgeEndStarCommon) computeLabelling(geomGraph []GeometryGraph) {
	ees.computeEdgeEndLabels(geomGraph[0].boundaryNodeRule)
	// Propagate side labels  around the edges in the star
	// for each parent Geometry
	ees.propagateSideLabels(0)
	ees.propagateSideLabels(1)

	// If there are edges that still have null labels for a geometry
	// this must be because there are no area edges for that geometry incident on this node.
	// In this case, to label the edge for that geometry we must test whether the
	// edge is in the interior of the geometry.
	// To do this it suffices to determine whether the node for the edge is in the interior of an area.
	// If so, the edge has location INTERIOR for the geometry.
	// In all other cases (e.g. the node is on a line, on a point, or not on the geometry at all) the edge
	// has the location EXTERIOR for the geometry.
	//
	// Note that the edge cannot be on the BOUNDARY of the geometry, since then
	// there would have been a parallel edge from the Geometry at this node also labelled BOUNDARY
	// and this edge would have been labelled in the previous step.
	//
	// This code causes a problem when dimensional collapses are present, since it may try and
	// determine the location of a node where a dimensional collapse has occurred.
	// The point should be considered to be on the EXTERIOR
	// of the polygon, but locate() will return INTERIOR, since it is passed
	// the original Geometry, not the collapsed version.
	//
	// If there are incident edges which are Line edges labelled BOUNDARY,
	// then they must be edges resulting from dimensional collapses.
	// In this case the other edges can be labelled EXTERIOR for this Geometry.
	//
	// MD 8/11/01 - NOT TRUE!  The collapsed edges may in fact be in the interior of the Geometry,
	// which means the other edges should be labelled INTERIOR for this Geometry.
	// Not sure how solve this...  Possibly labelling needs to be split into several phases:
	// area label propagation, symLabel merging, then finally null label resolution.
	hasDimensionalCollapseEdge := []bool{false, false}

	ees.edgeMap.Walk(func(key, value interface{}) {
		e := value.(edgeEnd)
		label := e.getLabel()
		for geomi := 0; geomi < 2; geomi++ {
			if label[geomi].isLine() && label[geomi][ON] == location.Boundary {
				hasDimensionalCollapseEdge[geomi] = true
			}
		}
	})

	ees.edgeMap.Walk(func(key, value interface{}) {
		e := value.(edgeEnd)
		label := e.getLabel()

		for geomi := 0; geomi < 2; geomi++ {
			if label[geomi].isNull() {
				loc := location.None
				if hasDimensionalCollapseEdge[geomi] {
					loc = location.Exterior
				} else {
					p := e.getCoord()
					loc = ees.getLocation(geomi, p, geomGraph)
				}
				label[geomi].setAllLocationsIfNull(loc)
			}
		}
	})
}
func (ees *edgeEndStarCommon) computeEdgeEndLabels(boundaryNodeRule boundary.NodeRule) {
	// Compute edge label for each EdgeEnd
	ees.edgeMap.Walk(func(key, value interface{}) {
		ee := value.(edgeEnd)
		ee.computeLabel(boundaryNodeRule)
	})
}

func (ees *edgeEndStarCommon) getLocation(geomIndex int, p geom.Coord, geom []GeometryGraph) location.Type {
	// compute location only on demand
	if ees.ptInAreaLocation[geomIndex] == location.None {
		ees.ptInAreaLocation[geomIndex] = xy.LocatePointOnGeomSFSBoundaryRun(p, geom[geomIndex].parentGeom)
	}
	return ees.ptInAreaLocation[geomIndex]
}
func (ees *edgeEndStarCommon) isAreaLabelsConsistent(geomGraph *GeometryGraph) (bool, error) {
	ees.computeEdgeEndLabels(geomGraph.boundaryNodeRule)
	return ees.checkAreaLabelsConsistent(0)
}

func (ees *edgeEndStarCommon) checkAreaLabelsConsistent(geomIndex int) (bool, error) {
	// Since edges are stored in CCW order around the node,
	// As we move around the ring we move from the right to the left side of the edge

	// if no edges, trivially consistent
	if ees.edgeMap.Size() <= 0 {
		return true, nil
	}

	// initialize startLoc to location of last L side (if any)
	lastEdgeIndex := ees.edgeMap.Size() - 1
	startLabel := ees.getEdgeEnd(lastEdgeIndex).getLabel()
	startLoc := startLabel[geomIndex][LEFT]

	if startLoc == location.None {
		return false, fmt.Errorf("Found unlabelled area edge")
	}

	currLoc := startLoc
	found := true

	var err error = nil
	ees.edgeMap.WalkInterruptible(func(key, value interface{}) bool {
		e := value.(edgeEnd)
		label := e.getLabel()
		// we assume that we are only checking a area
		if label[geomIndex].isArea() {
			err = fmt.Errorf("Found non-area edge")
		}

		leftLoc := label[geomIndex][LEFT]
		rightLoc := label[geomIndex][RIGHT]

		// check that edge is really a boundary between inside and outside!
		if leftLoc == rightLoc {
			found = false
			return false
		}
		// check side location conflict
		if rightLoc != currLoc {
			found = false
			return false
		}
		currLoc = leftLoc
		return true
	})

	return found, err
}
func (ees *edgeEndStarCommon) propagateSideLabels(geomIndex int) error {
	// Since edges are stored in CCW order around the node,
	// As we move around the ring we move from the right to the left side of the edge
	startLoc := location.None

	// initialize loc to location of last L side (if any)
	//System.out.println("finding start location");
	ees.edgeMap.Walk(func(key, value interface{}) {
		e := value.(edgeEnd)
		label := e.getLabel()
		if label[geomIndex].isArea() && label[geomIndex][LEFT] != location.None {
			startLoc = label[geomIndex][LEFT]
		}
	})

	// no labelled sides found, so no labels to propagate
	if startLoc == location.None {
		return nil
	}

	currLoc := startLoc
	var err error = nil
	ees.edgeMap.Walk(func(key, value interface{}) {
		e := value.(edgeEnd)
		label := e.getLabel()
		// set null ON values to be in current location
		if label[geomIndex][ON] == location.None {
			label[geomIndex][ON] = currLoc
		}
		// set side labels (if any)
		if label[geomIndex].isArea() {
			leftLoc := label[geomIndex][LEFT]
			rightLoc := label[geomIndex][RIGHT]
			// if there is a right location, that is the next location to propagate
			if rightLoc != location.None {
				if rightLoc != currLoc {
					err = fmt.Errorf("side location conflict %v", e.getCoord())
				}
			}
			if leftLoc == location.None {
				err = fmt.Errorf("found single null side (at %v)", e.getCoord())
			}
			currLoc = leftLoc
		} else {
			/** RHS is null - LHS must be null too.
			 *  This must be an edge from the other geometry, which has no location
			 *  labelling for this geometry.  This edge must lie wholly inside or outside
			 *  the other geometry (which is determined by the current location).
			 *  Assign both sides to be the current location.
			 */
			if label[geomIndex][LEFT] != location.None {
				err = fmt.Errorf("found single null side")
			}
			label[geomIndex][RIGHT] = currLoc
			label[geomIndex][LEFT] = currLoc
		}
	})

	return err
}

func (ees *edgeEndStarCommon) findIndex(eSearch edgeEnd) int {
	found := false
	i := 0
	ees.edgeMap.WalkInterruptible(func(key, value interface{}) bool {
		i++
		if value == eSearch {
			found = true
			return false
		}
		return true
	})

	if found {
		return i
	}

	return -1
}

func (ees *edgeEndStarCommon) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("EdgeEndStar: %v\n", ees.coordinate()))

	ees.edgeMap.Walk(func(key, e interface{}) {
		buf.WriteString(fmt.Sprintf("%v\n", e))
	})
	return buf.String()
}
