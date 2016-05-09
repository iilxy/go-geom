package xy_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"math"
	"reflect"
	"testing"
)

func TestPointCentroidCalculator_GetCentroid_NoCoordsAdded(t *testing.T) {
	calculator := xy.NewPointCentroidCalculator()
	centroid := calculator.GetCentroid()
	if !centroid.Equal(geom.XY, geom.Coord{math.NaN(), math.NaN()}) {
		t.Errorf("centroid with no coords added should return the [NaN NaN] coord but was: %v", centroid)
	}
}

type pointTestData struct {
	points []*geom.Point
	result geom.Coord
}

func TestPointGetCentroid(t *testing.T) {
	for i, tc := range []pointTestData{
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 2}),
			},
			result: geom.Coord{1, 1},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 0}),
			},
			result: geom.Coord{1, 0},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 2}),
				geom.NewPointFlat(geom.XY, []float64{0, 2}),
			},
			result: geom.Coord{1, 1},
		},
	} {
		checkPointsCentroidFunc(t, i, tc)
		checkPointCentroidFlatFunc(t, i, tc)
		checkPointCentroidMultiPoint(t, i, tc)
		checkAddEachPoint(t, i, tc)

	}

}

func checkPointsCentroidFunc(t *testing.T, i int, tc pointTestData) {
	centroid := xy.PointsCentroid(tc.points[0], tc.points[1:]...)

	if !reflect.DeepEqual(tc.result, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.result, centroid)
	}

}
func checkPointCentroidFlatFunc(t *testing.T, i int, tc pointTestData) {
	data := make([]float64, len(tc.points)*2, len(tc.points)*2)

	for i, p := range tc.points {
		data[i*2] = p.FlatCoords()[0]
		data[(i*2)+1] = p.FlatCoords()[1]
	}
	centroid := xy.PointsCentroidFlat(geom.XY, data)

	if !reflect.DeepEqual(tc.result, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.result, centroid)
	}

}
func checkPointCentroidMultiPoint(t *testing.T, i int, tc pointTestData) {
	centroid := xy.MultiPointCentroid(createMultiPoint(tc.points))

	if !reflect.DeepEqual(tc.result, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for multipoint to be\n%v but was \n%v", i+1, tc.result, centroid)
	}
}

func createMultiPoint(points []*geom.Point) *geom.MultiPoint {
	data := make([]float64, len(points)*2, len(points)*2)

	for i, p := range points {
		flatCoords := p.FlatCoords()
		data[i*2] = flatCoords[0]
		data[(i*2)+1] = flatCoords[1]
	}
	return geom.NewMultiPointFlat(geom.XY, data)
}

func checkAddEachPoint(t *testing.T, i int, tc pointTestData) {
	calc := xy.NewPointCentroidCalculator()
	for _, p := range tc.points {
		calc.AddPoint(p)
	}
	centroid := calc.GetCentroid()

	if !reflect.DeepEqual(tc.result, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.result, centroid)
	}

}
