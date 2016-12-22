package geom

import (
	"github.com/twpayne/go-geom/internal/tree"
	"github.com/twpayne/go-geom/xy/boundary"
)

func boundary1(flatCoords []float64, layout Layout) T {
	if len(flatCoords) == 0 {
		return NewMultiPoint(layout)
	}

	bdyPts := computeBoundaryCoordinates(boundary.Mod2BoundaryNodeRule{}, flatCoords, layout, []int{0, len(flatCoords)})

	// return Point or MultiPoint
	if len(bdyPts) == 1 {
		return NewPointFlat(layout, bdyPts)
	}

	// this handles 0 points case as well
	return NewMultiPointFlat(layout, bdyPts)
}

func computeBoundaryCoordinates(boundaryRule boundary.NodeRule, flatCoords []float64, layout Layout, ends []int) []float64 {
	bdyPts := []float64{}
	endpointMap := tree.NewTreeMap(coordCompare{})
	endpointMap.SetDefault(0)
	stride := layout.Stride()

	start := 0
	for _, end := range ends {
		if start != end {
			addEndpoint(endpointMap, flatCoords[start:start+stride])
			addEndpoint(endpointMap, flatCoords[end-stride:end])
		}

		start = end
	}

	endpointMap.Walk(func(key, value interface{}) {
		valence := value.(int)
		if boundaryRule.IsInBoundary(valence) {
			coord := key.(Coord)
			bdyPts = append(bdyPts, coord...)
		}
	})

	return bdyPts
}

func addEndpoint(endpointMap *tree.TreeMap, coord []float64) {
	valence, _ := endpointMap.Get(Coord(coord))

	endpointMap.Insert(Coord(coord), valence.(int)+1)
}

type coordCompare struct{}

func (cc coordCompare) IsEquals(o1, o2 interface{}) bool {
	v1 := o1.(Coord)
	v2 := o2.(Coord)
	v1len := len(v1)

	if v1len != len(v2) {
		return false
	}

	for i := 0; i < v1len; i++ {
		if v1[i] != v2[i] {
			return false
		}
	}

	return true
}

func (cc coordCompare) IsLess(o1, o2 interface{}) bool {
	v1 := o1.(Coord)
	v2 := o2.(Coord)
	if v1[0] < v2[0] {
		return true
	}
	if v1[0] > v2[0] {
		return false
	}
	if v1[1] < v2[1] {
		return true
	}

	return false
}
