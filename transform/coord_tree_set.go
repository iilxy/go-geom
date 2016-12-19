package transform

import (
	"fmt"
	"github.com/twpayne/go-geom"
)

// Compare compares two coordinates for equality and magnitude
type CoordCompare interface {
	IsEquals(x, y geom.Coord) bool
	IsLess(x, y geom.Coord) bool
}

type compareAdapter struct {
	compare CoordCompare
}

var _ Compare = compareAdapter{}

func (c compareAdapter) IsEquals(o1, o2 interface{}) bool {
	return c.compare.IsEquals(o1.(geom.Coord), o2.(geom.Coord))
}
func (c compareAdapter) IsLess(o1, o2 interface{}) bool {
	return c.compare.IsLess(o1.(geom.Coord), o2.(geom.Coord))
}

// CoordTreeSet sorts the coordinates according to the Compare strategy and removes duplicates as
// dictated by the Equals function of the Compare strategy
type CoordTreeSet struct {
	treeMap *TreeMap
	layout  geom.Layout
	stride  int
}

// NewCoordTreeSet creates a new TreeSet instance
func NewCoordTreeSet(layout geom.Layout, compare CoordCompare) *CoordTreeSet {
	treeMap := NewTreeMap(compareAdapter{compare})
	return &CoordTreeSet{
		layout:  layout,
		stride:  layout.Stride(),
		treeMap: treeMap,
	}
}

// Insert adds a new coordinate to the tree set
// the coordinate must be the same size as the Stride of the layout provided
// when constructing the TreeSet
// Returns true if the coordinate was added, false if it was already in the tree
func (set *CoordTreeSet) Insert(coord geom.Coord) bool {
	if set.stride == 0 {
		set.stride = set.layout.Stride()
	}
	if len(coord) < set.stride {
		panic(fmt.Sprintf("Coordinate inserted into tree does not have a sufficient number of points for the provided layout.  Length of Coord was %v but should have been %v", len(coord), set.stride))
	}
	return set.treeMap.Insert(coord, nil)
}

// ToFlatArray returns an array of floats containing all the coordinates in the TreeSet
func (set *CoordTreeSet) ToFlatArray() []float64 {
	stride := set.layout.Stride()
	size := set.treeMap.Size()
	array := make([]float64, size*stride, size*stride)

	i := 0
	set.treeMap.Walk(func(k, v interface{}) {
		coord := k.(geom.Coord)
		for j := 0; j < stride; j++ {
			array[i+j] = coord[j]
		}
		i += stride
	})

	return array
}
