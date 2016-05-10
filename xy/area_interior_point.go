package xy

import "github.com/twpayne/go-geom"

// PolygonsInteriorPoint computes a point in the interior of one of the provided polygons
//
// Algorithm
// * Find the intersections between the geometry and the horizontal bisector of the area's envelope
// * Pick the midpoint of the largest intersection (the intersections will be lines and points)
func PolygonsInteriorPoint(polygon *geom.Polygon, extraPolys ...*geom.Polygon) (centroid geom.Coord) {

}

// PolygonsInteriorPoint computes a point in the interior of one of the polygons in the MultiPolygon
//
// Algorithm
// * Find the intersections between the geometry and the horizontal bisector of the area's envelope
// * Pick the midpoint of the largest intersection (the intersections will be lines and points)
func MultiPolygonInteriorPoint(polygon *geom.MultiPolygon) (centroid geom.Coord) {

}

type areaInteriorPointCalculator struct {
	interiorPoint geom.Coord
	maxWidth      float64
}

func (calc *areaInteriorPointCalculator) addPolygon(geometry *geom.Polygon) {
	bisector := calc.horizontalBisector(geometry)

	intersections := bisector.intersection(geometry)
	widestIntersection := calc.widestGeometry(intersections)

	width := widestIntersection.Bounds().Length(0)
	if len(calc.interiorPoint) == 0 || width > calc.maxWidth {
		calc.interiorPoint = widestIntersection.Bounds().Center()
		calc.maxWidth = width
	}
}
func (calc *areaInteriorPointCalculator) widestGeometry(gc []geom.T) geom.T {
	if len(gc) == 0 {
		return nil
	}

	widestGeometry := gc[0]
	for i := 1; i < len(gc); i++ {
		//Start at 1
		if gc[i].Bounds().Length(0) > widestGeometry.Bounds().Length(0) {
			widestGeometry = gc[i]
		}
	}
	return widestGeometry
}

func (calc *areaInteriorPointCalculator) horizontalBisector(geometry geom.T) geom.LineString {
	envelope := geometry.Bounds()

	// Assert: for areas, minx <> maxx
	avgY := avg(envelope.Min(1), envelope.Max(1))
	return geom.NewLineStringFlat(geom.XY, []float64{envelope.Min(0), avgY, envelope.Max(0), avgY})
}

func avg(a, b float64) float64 {
	return (a + b) / 2.0
}
