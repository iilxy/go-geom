package lineintersector

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/lineintersection"
	"math"
)

// IsOnLine tests whether a point lies on the line segments defined by a list of
// coordinates.
//
// Returns true if the point is a vertex of the line or lies in the interior
//         of a line segment in the linestring
func IsOnLine(layout geom.Layout, point geom.Coord, lineSegmentCoordinates []float64) bool {
	stride := layout.Stride()
	if len(lineSegmentCoordinates) < (2 * stride) {
		panic(fmt.Sprintf("At least two coordinates are required in the lineSegmentsCoordinates array in 'algorithms.IsOnLine', was: %v", lineSegmentCoordinates))
	}
	strategy := RobustLineIntersector{}

	for i := stride; i < len(lineSegmentCoordinates); i += stride {
		segmentStart := lineSegmentCoordinates[i-stride : i-stride+2]
		segmentEnd := lineSegmentCoordinates[i : i+2]

		if PointIntersectsLine(strategy, geom.Coord(point), geom.Coord(segmentStart), geom.Coord(segmentEnd)) {
			return true
		}
	}
	return false
}

// Strategy is the line intersection implementation
type Strategy interface {
	computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord)
	computeLineOnLineIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord)
}

// PointIntersectsLine tests if point intersects the line
func PointIntersectsLine(strategy Strategy, point, lineStart, lineEnd geom.Coord) (hasIntersection bool) {
	intersectorData := &lineIntersectorData{
		strategy:           strategy,
		inputLines:         [2][2]geom.Coord{[2]geom.Coord{lineStart, lineEnd}, [2]geom.Coord{}},
		intersectionPoints: [2]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	strategy.computePointOnLineIntersection(intersectorData, point, lineStart, lineEnd)

	return intersectorData.intersectionType != lineintersection.NoIntersection
}

// LineIntersectsLine tests if the first line (line1Start,line1End) intersects the second line (line2Start, line2End)
// and returns a data structure that indicates if there was an intersection, the type of intersection and where the intersection
// was.  See lineintersection.Result for a more detailed explanation of the result object
func LineIntersectsLine(strategy Strategy, line1Start, line1End, line2Start, line2End geom.Coord) lineintersection.Result {
	intersectorData := &lineIntersectorData{
		strategy:           strategy,
		inputLines:         [2][2]geom.Coord{[2]geom.Coord{line2Start, line2End}, [2]geom.Coord{line1Start, line1End}},
		intersectionPoints: [2]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	strategy.computeLineOnLineIntersection(intersectorData, line1Start, line1End, line2Start, line2End)

	var intersections []geom.Coord

	switch intersectorData.intersectionType {
	case lineintersection.NoIntersection:
		intersections = []geom.Coord{}
	case lineintersection.PointIntersection:
		intersections = intersectorData.intersectionPoints[:1]
	case lineintersection.CollinearIntersection:
		intersections = intersectorData.intersectionPoints[:2]
	}
	return lineintersection.NewResult(intersectorData.intersectionType, intersections)
}

// An internal data structure for containing the data during calculations
type lineIntersectorData struct {
	indexComputed bool
	// new Coordinate[2][2];
	inputLines [2][2]geom.Coord

	// if only a point intersection then 0 index coord will contain the intersection point
	// if co-linear (lines overlay each other) the two coordinates represent the start and end points of the overlapping lines.
	intersectionPoints [2]geom.Coord
	intersectionType   lineintersection.Type

	// The indexes of the endpoints of the intersection lines, in order along
	// the corresponding line
	intLineIndex [2][2]int
	isProper     bool
	pa, pb       geom.Coord
	strategy     Strategy
}

/**
 *  RParameter computes the parameter for the point p
 *  in the parameterized equation
 *  of the line from p1 to p2.
 *  This is equal to the 'distance' of p along p1-p2
 */
func rParameter(p1, p2, p geom.Coord) float64 {
	var r float64
	// compute maximum delta, for numerical stability
	// also handle case of p1-p2 being vertical or horizontal
	dx := math.Abs(p2[0] - p1[0])
	dy := math.Abs(p2[1] - p1[1])
	if dx > dy {
		r = (p[0] - p1[0]) / (p2[0] - p1[0])
	} else {
		r = (p[1] - p1[1]) / (p2[1] - p1[1])
	}
	return r
}
