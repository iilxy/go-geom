package relate

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/internal"
	"github.com/twpayne/go-geom/xy/location"
)

type Calculation struct {
	dimGeomA, dimGeomB int
	matrix             internal.IntersectionMatrix
}

func (rc Calculation) Within() bool {
	return rc.matrix.Within()
}

func (rc Calculation) Contains() bool {
	return rc.matrix.Contains()
}

func (rc Calculation) CoveredBy() bool {
	return rc.matrix.CoveredBy()
}

func (rc Calculation) Crosses() bool {
	return rc.matrix.Crosses(rc.dimGeomA, rc.dimGeomB)
}

func (rc Calculation) Covers() bool {
	return rc.matrix.Covers()
}

func (rc Calculation) Disjoint() bool {
	return rc.matrix.Disjoint()
}

func (rc Calculation) Equal() bool {
	return rc.matrix.Equal(rc.dimGeomA, rc.dimGeomB)
}

func (rc Calculation) Overlaps() bool {
	return rc.matrix.Overlaps(rc.dimGeomA, rc.dimGeomB)
}

func (rc Calculation) Touches() bool {
	return rc.matrix.Touches(rc.dimGeomA, rc.dimGeomB)
}

func Calculate(geom1, geom2 geom.T) Calculation {
	im := internal.IntersectionMatrix{}

	// since Geometries are finite and embedded in a 2-D space, the EE element must always be 2
	im[location.Exterior][location.Exterior] = 2

	// if the Geometries don't overlap there is nothing to do
	if !geom1.Bounds().Overlaps(geom1.Layout(), geom2.Bounds()) {
		computeDisjointIM(im, geom1, geom2)
		return im
	}
}

func computeDisjointIM(im internal.IntersectionMatrix, geom1, geom2 geom.T) {
	if len(geom1.FlatCoords()) == 0 {

		im[location.Interior][location.Interior] = geom1.Layout().Stride()
		im[location.Boundary][location.Exterior] = ga.getBoundaryDimension()
	}
}
