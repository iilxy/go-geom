package xy_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/internal"
	"reflect"
	"testing"
)

func TestPointInteriorPoint(t *testing.T) {
	for i, tc := range []pointTestData{
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 2}),
			},
			result: geom.Coord{0, 0},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{1.5, 1.5}),
				geom.NewPointFlat(geom.XY, []float64{2, 2}),
			},
			result: geom.Coord{1.5, 1.5},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 0}),
				geom.NewPointFlat(geom.XY, []float64{1, 0}),
			},
			result: geom.Coord{1, 0},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
			},
			result: geom.Coord{0, 0},
		},
	} {
		interiorPoint := xy.PointsInteriorPoint(tc.points[0], tc.points[1:]...)
		if !reflect.DeepEqual(interiorPoint, tc.result) {
			t.Errorf("Test %v failed: xy.PointsInteriorPoint(tc.points[0], tc.points[1:]...). Expected \n\t%v but was \n\t%v", i+1, tc.result, interiorPoint)
		}

		interiorPoint = xy.MultiPointInteriorPoint(createMultiPoint(tc.points))
		if !reflect.DeepEqual(interiorPoint, tc.result) {
			t.Errorf("Test %v failed:xy.MultiPointInteriorPoint(createMultiPoint(tc.points). Expected \n\t%v but was \n\t%v", i+1, tc.result, interiorPoint)
		}
	}
}

func TestPointsInteriorPoint(t *testing.T) {
	multiPoint := geom.NewMultiPointFlat(geom.XY, internal.RING.FlatCoords())
	interiorPoint := xy.MultiPointInteriorPoint(multiPoint)
	expected := geom.Coord{-71.1019285062273, 42.3147384934248}

	if !reflect.DeepEqual(interiorPoint, expected) {
		t.Errorf("Expected: \n\t%v but was\n\t%v", expected, interiorPoint)
	}
}
