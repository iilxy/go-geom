package xy

import (
	"github.com/twpayne/go-geom"
	"math"
)

// PointInteriorPoint returns the point closest to the centroid of all the points
func PointsInteriorPoint(point *geom.Point, extra ...*geom.Point) geom.Coord {
	centroid := PointsCentroid(point, extra...)

	calc := newPointInteriorPointCalc(centroid)
	calc.add(point.Layout(), point.FlatCoords())

	for _, p := range extra {
		calc.add(p.Layout(), p.FlatCoords())
	}

	return calc.interiorPoint
}

// PointInteriorPoint returns the point closest to the centroid of all the points in the multi-point
func MultiPointInteriorPoint(point *geom.MultiPoint) geom.Coord {

	centroid := MultiPointCentroid(point)
	calc := newPointInteriorPointCalc(centroid)
	calc.add(point.Layout(), point.FlatCoords())
	return calc.interiorPoint

}

type pointInteriorPointCalc struct {
	centroid, interiorPoint geom.Coord
	minDistance             float64
}

func newPointInteriorPointCalc(centroid geom.Coord) *pointInteriorPointCalc {
	return &pointInteriorPointCalc{
		minDistance: math.MaxFloat64,
		centroid:    centroid,
	}
}

func (calc *pointInteriorPointCalc) add(layout geom.Layout, coords []float64) {
	stride := layout.Stride()
	for i := 0; i < len(coords); i += stride {
		point := geom.Coord(coords[i : i+stride])
		dist := Distance(point, calc.centroid)
		if dist < calc.minDistance {
			calc.interiorPoint = make(geom.Coord, stride)
			copy(calc.interiorPoint, point)
			calc.minDistance = dist
		}
	}
}
