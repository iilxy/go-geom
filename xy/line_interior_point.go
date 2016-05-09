package xy

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"math"
)

// LinesInteriorPoint computes the interior point of all the LineStrings provided as arguments.
//
// Algorithm:
// * Find an interior vertex which is closest to the centroid of the linestring.
// * If there is no interior vertex, find the endpoint which is closest to the centroid.
func LinesInteriorPoint(line *geom.LineString, extraLines ...*geom.LineString) (interiorPoint geom.Coord) {
	centroid := LinesCentroid(line, extraLines...)
	fmt.Println(centroid)
	calc := newLineInteriorPointCalc(centroid)
	calc.addInteriorPoints(line.Layout(), line.FlatCoords())
	for _, extraLine := range extraLines {
		calc.addInteriorPoints(extraLine.Layout(), extraLine.FlatCoords())
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

func (calc *lineInteriorPointCalc) addInteriorPoints(layout geom.Layout, coords []float64) (interiorPointWasSet bool) {
	stride := layout.Stride()
	interiorPointWasSet = false
	for i := 0; i < len(coords); i += stride {
		setNow := calc.addCoord(geom.Coord(coords[i : i+stride]))
		interiorPointWasSet = interiorPointWasSet || setNow

	}

	return interiorPointWasSet
}

func (calc *lineInteriorPointCalc) addCoord(point geom.Coord) bool {
	dist := Distance(point, calc.centroid)
	if dist < calc.minDistance {
		calc.interiorPoint = make(geom.Coord, len(point))
		copy(calc.interiorPoint, point)
		calc.minDistance = dist
		return true
	}

	fmt.Println(calc.minDistance, dist)
	fmt.Println(point, calc.minDistance)
	fmt.Println()

	return false
}
