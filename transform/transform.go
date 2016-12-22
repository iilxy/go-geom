package transform

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/internal"
)

// UniqueCoords creates a new coordinate array (with the same layout as the inputs) that
// contains each unique coordinate in the coordData.  The ordering of the coords are the
// same as the input.
func UniqueCoords(layout geom.Layout, compare internal.CoordCompare, coordData []float64) []float64 {
	set := internal.NewCoordTreeSet(layout, compare)
	stride := layout.Stride()
	uniqueCoords := make([]float64, 0, len(coordData))
	numCoordsAdded := 0
	for i := 0; i < len(coordData); i += stride {
		coord := coordData[i : i+stride]
		added := set.Insert(geom.Coord(coord))

		if added {
			uniqueCoords = append(uniqueCoords, coord...)
			numCoordsAdded++
		}
	}
	return uniqueCoords[:numCoordsAdded*stride]
}
