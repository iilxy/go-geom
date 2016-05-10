package xy

import (
	"math"

	"github.com/twpayne/go-geom"
)

// LinesInteriorPoint computes the interior point of all the LineStrings provided as arguments.
//
// Algorithm:
// * Find an interior vertex which is closest to the centroid of the linestring.
// * If there is no interior vertex, find the endpoint which is closest to the centroid.
func LinesInteriorPoint(line *geom.LineString, extraLines ...*geom.LineString) (interiorPoint geom.Coord) {
	centroid := LinesCentroid(line, extraLines...)
	calc := newLineInteriorPointCalc(centroid)
	calc.addInteriorPoints(line.Layout(), line.FlatCoords())
	for _, extraLine := range extraLines {
		calc.addInteriorPoints(extraLine.Layout(), extraLine.FlatCoords())
	}
	if len(calc.interiorPoint) == 0 {
		calc.addEndpoints(line.Layout(), line.FlatCoords())
		for _, extraLine := range extraLines {
			calc.addEndpoints(extraLine.Layout(), extraLine.FlatCoords())
		}
	}
	return calc.interiorPoint
}

// LineRingInteriorPoint computes the interior point of all the LinearRings provided as arguments.
//
// Algorithm:
// * Find an interior vertex which is closest to the centroid of the linestring.
// * If there is no interior vertex, find the endpoint which is closest to the centroid.
func LinearRingsInteriorPoint(line *geom.LinearRing, extraLines ...*geom.LinearRing) (interiorPoint geom.Coord) {
	centroid := LinearRingsCentroid(line, extraLines...)
	calc := newLineInteriorPointCalc(centroid)
	calc.addInteriorPoints(line.Layout(), line.FlatCoords())
	for _, extraLine := range extraLines {
		calc.addInteriorPoints(extraLine.Layout(), extraLine.FlatCoords())
	}

	if len(calc.interiorPoint) == 0 {
		calc.addEndpoints(line.Layout(), line.FlatCoords())
		for _, extraLine := range extraLines {
			calc.addEndpoints(extraLine.Layout(), extraLine.FlatCoords())
		}
	}

	return calc.interiorPoint
}

// MultiLineInteriorPoint computes the interior point of the MultiLineString string
//
// Algorithm:
// * Find an interior vertex which is closest to the centroid of the linestring.
// * If there is no interior vertex, find the endpoint which is closest to the centroid.
func MultiLineInteriorPoint(line *geom.MultiLineString) (interiorPoint geom.Coord) {
	centroid := MultiLineCentroid(line)
	calc := newLineInteriorPointCalc(centroid)
	calc.addInteriorPoints(line.Layout(), line.FlatCoords())

	if len(calc.interiorPoint) == 0 {
		layout := line.Layout()
		for i, n := 0, line.NumLineStrings(); i < n; i++ {
			calc.addEndpoints(layout, line.LineString(i).FlatCoords())
		}
	}
	return calc.interiorPoint
}

type lineInteriorPointCalc struct {
	centroid, interiorPoint geom.Coord
	minDistance             float64
}

func newLineInteriorPointCalc(centroid geom.Coord) *lineInteriorPointCalc {
	return &lineInteriorPointCalc{
		minDistance: math.MaxFloat64,
		centroid:    centroid,
	}
}

func (calc *lineInteriorPointCalc) addInteriorPoints(layout geom.Layout, coords []float64) {
	stride := layout.Stride()
	for i := stride; i < len(coords)-stride; i += stride {
		calc.addCoord(geom.Coord(coords[i : i+stride]))
	}
}

func (calc *lineInteriorPointCalc) addCoord(point geom.Coord) {
	dist := Distance(point, calc.centroid)
	if dist < calc.minDistance {

		switch {
		case len(calc.interiorPoint) < len(point):
			calc.interiorPoint = make(geom.Coord, len(point))
		case len(calc.interiorPoint) > len(point):
			calc.interiorPoint = calc.interiorPoint[:len(point)]
		}

		copy(calc.interiorPoint, point)
		calc.minDistance = dist
	}
}

func (calc *lineInteriorPointCalc) addEndpoints(layout geom.Layout, coords []float64) {
	calc.addCoord(geom.Coord(coords[:layout.Stride()]))
	calc.addCoord(geom.Coord(coords[len(coords)-layout.Stride():]))
}
