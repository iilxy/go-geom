package graph

import (
	"fmt"
	"github.com/twpayne/go-geom/xy/location"
)

const nullDepth = -1

// Depth objects record the topological depth of the sides
// of an Edge for up to two Geometries.
type Depth [2][3]int

func newDepth() Depth {
	return Depth{
		[3]int{nullDepth, nullDepth, nullDepth},
		[3]int{nullDepth, nullDepth, nullDepth},
	}
}

func (d *Depth) location(geomIndex, posIndex int) location.Type {
	if d[geomIndex][posIndex] <= 0 {
		return location.Exterior
	}
	return location.Interior
}

func (d *Depth) addLocation(geomIndex, posIndex, loc location.Type) {
	if loc == location.Interior {
		d[geomIndex][posIndex]++
	}
}

func (d *Depth) addFromLabel(lbl Label) {
	for i := 0; i < 2; i++ {
		for j := 1; j < 3; j++ {
			loc := lbl[i][j]
			if loc == location.Exterior || loc == location.Interior {
				// initialize depth if it is null, otherwise add this location value
				if d.isNullPos(i, j) {
					d[i][j] = depthAtLocation(loc)
				} else {
					d[i][j] += depthAtLocation(loc)
				}
			}
		}
	}
}

func (d *Depth) delta(geomIndex int) int {
	return d[geomIndex][RightOfLabel] - d[geomIndex][LeftOfLabel]
}

// isNull returns if all depths are null (uninitialized)
func (d *Depth) isNull() bool {
	for _, geomIdx := range d {
		for _, pos := range geomIdx {
			if pos != nullDepth {
				return false
			}
		}
	}

	return true
}

func (d *Depth) isNullGeom(geomIndex int) bool {
	return d[geomIndex][1] == nullDepth
}

func (d *Depth) isNullPos(geomIndex, posIndex int) bool {
	return d[geomIndex][posIndex] == nullDepth
}

// Normalize the depths for each geometry, if they are non-null.
// A normalized depth has depth values in the set { 0, 1 }.
// Normalizing the depths involves reducing the depths by the same amount so that at least
// one of them is 0.  If the remaining value is > 0, it is set to 1.
func (d *Depth) normalize() {
	for i := 0; i < 2; i++ {
		if !d.isNullGeom(i) {
			minDepth := d[i][1]
			if d[i][2] < minDepth {
				minDepth = d[i][2]

				if minDepth < 0 {
					minDepth = 0
				}
				for j := 1; j < 3; j++ {
					newValue := 0
					if d[i][j] > minDepth {
						newValue = 1
					}
					d[i][j] = newValue
				}
			}
		}
	}
}

func (d *Depth) String() string {
	return fmt.Sprintf("A: %v,%v B: %v,%v", d[0][1], d[0][2], d[1][1], d[1][2])
}

func depthAtLocation(loc location.Type) int {
	if loc == location.Exterior {
		return 0
	}
	if loc == location.Interior {
		return 1
	}
	return nullDepth
}
