package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom"
)

type quadrant int

// Utility functions for working with quadrants, which are numbered as follows:
//
//     1 | 0
//     --+--
//     2 | 3
//
const (
	NE quadrant = iota
	NW
	SW
	SE
)

//doublesQuadrant returns the quadrant of a directed line segment (specified as x and y
// displacements, which cannot both be 0).
func doublesQuadrant(dx, dy float64) quadrant {
	if dx == 0.0 && dy == 0.0 {
		panic(fmt.Sprintf("Cannot compute the quadrant for point (%v, %v)", dx, dy))
	}

	if dx >= 0.0 {
		if dy >= 0.0 {
			return NE
		}
		return SE
	}

	if dy >= 0.0 {
		return NW
	}
	return SW
}

// coordsQuadrant returns the quadrant of a directed line segment from p0 to p1.
func coordsQuadrant(p0, p1 geom.Coord) quadrant {
	if p1[0] == p0[0] && p1[1] == p0[1] {
		panic(fmt.Sprintf("Cannot compute the quadrant for two identical points %v", p0))
	}

	if p1[0] >= p0[0] {
		if p1[1] >= p0[1] {
			return NE
		}
		return SE
	}

	if p1[1] >= p0[1] {
		return NW
	}
	return SW
}

// isOpposite returns true if the quadrants are 1 and 3, or 2 and 4
func (quad1 quadrant) isOpposite(quad2 quadrant) bool {
	if quad1 == quad2 {
		return false
	}

	diff := (quad1 - quad2 + 4) % 4

	// if quadrants are not adjacent, they are opposite
	if diff == 2 {
		return true
	}
	return false
}

// commonHalfPlane returns the right-hand quadrant of the halfplane defined by the two quadrants,
// or -1 if the quadrants are opposite, or the quadrant if they are identical.
func (quad1 quadrant) commonHalfPlane(quad2 quadrant) quadrant {
	// if quadrants are the same they do not determine a unique common halfplane.
	// Simply return one of the two possibilities
	if quad1 == quad2 {
		return quad1
	}
	diff := (quad1 - quad2 + 4) % 4
	// if quadrants are not adjacent, they do not share a common halfplane
	if diff == 2 {
		return -1
	}

	min := quad2
	if quad1 < quad2 {
		min = quad1
	}

	max := quad2
	if quad1 > quad2 {
		max = quad1
	}

	// for this one case, the righthand plane is NOT the minimum index;
	if min == 0 && max == 3 {
		return 3
	}

	// in general, the halfplane index is the minimum of the two adjacent quadrants
	return min
}

// isInHalfPlane Returns whether the given quadrant lies within the given halfplane (specified
// by its right-hand quadrant).
func (quad quadrant) isInHalfPlane(halfPlane quadrant) bool {

	if halfPlane == SE {
		return quad == SE || quad == SW
	}
	return quad == halfPlane || quad == halfPlane+1
}

func (quad quadrant) isNorthern() bool {
	return quad == NE || quad == NW
}
