package xy_test

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func ExamplePointsInteriorPoint() {
	interiorPoint := xy.PointsInteriorPoint(
		geom.NewPointFlat(geom.XY, []float64{0, 0}),
		geom.NewPointFlat(geom.XY, []float64{1.5, 1.5}),
		geom.NewPointFlat(geom.XY, []float64{2, 2}))
	fmt.Println(interiorPoint)
	// Output: [1.5 1.5]
}

func ExampleMultiPointInteriorPoint() {
	multiPoint := geom.NewMultiPointFlat(geom.XY, []float64{0, 0, 1.5, 1.5, 2, 2})
	interiorPoint := xy.MultiPointInteriorPoint(multiPoint)
	fmt.Println(interiorPoint)
	// Output: [1.5 1.5]

}
