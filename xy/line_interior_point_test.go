package xy_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"reflect"
	"testing"
)

func TestLinesInteriorPoint(t *testing.T) {
	for _, tc := range []lineDataType{
		{
			desc: "Closest to one endpoint",
			lines: []*geom.LineString{
				geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 10}),
				geom.NewLineStringFlat(geom.XY, []float64{0, 2, 10, 20}),
			},
			lineCentroid: geom.Coord{10, 10},
		},
		{
			desc: "Randomly Generated 1",
			lines: []*geom.LineString{
				geom.NewLineStringFlat(geom.XYZ, []float64{
					-2.0476962378424928E9, -1.1237663867526606E8, -4.4174369894724345E8, 1.6153951634676201E9, 1.9928170320125375E7, 6.047217789191747E8, 2.0108606440657005E7, -1.4446186281495008E9, 3.509100203058657E8, -2.0476962378424928E9, -1.1237663867526606E8, -4.4174369894724345E8,
				}),
			},
			lineCentroid: geom.Coord{2.0108606440657005E7, -1.4446186281495008E9, 3.509100203058657E8},
		},
		{
			desc: "Randomly Generated 2",
			lines: []*geom.LineString{
				geom.NewLineStringFlat(geom.XYZ, []float64{
					-7.46789843875317E8, 1.065439658580128E9, 7.816772072530246E8, 4.080231418111146E8, 1.2013259393533134E9, -1.2156916596965182E9, -1.074743579462955E8, -1.60593803795027E9, -1.8410942424943504E9, -2.0439198291142964E8, -2.0596380989404094E9, -1997897.6567936332, -1.5723919052723333E8, -2.6166424513434714E8, 8.704697990334049E7, -1.696389140551475E9, 4.492291336977326E8, -4.052487033168656E8, 1.1947427435988931E7, -7.092998017420655E8, -1.2269501402370174E8, -7.46789843875317E8, 1.065439658580128E9, 7.816772072530246E8,
				}),
			},
			lineCentroid: geom.Coord{-1.5723919052723333E8, -2.6166424513434714E8, 8.704697990334049E7},
		},
		{
			desc: "Randomly Generated 3",
			lines: []*geom.LineString{
				geom.NewLineStringFlat(geom.XYZ, []float64{
					6.341783697427113E7, -1.0726206418402787E8, -1.0839573118002522E9, 7.602495859103782E8, -5.431471681945347E8, -5.880492197023277E8, 2.628963870764161E8, 4.7620800643484056E8, 7.600627349557235E8, 2.7731211674290743E7, -1.5962198643569078E9, 1.1273405365563903E9, -2.72742709910209E8, 8.401246704749283E8, 8.139560428938211E8, -1.577217923210562E9, 5.0043679987498835E7, -1.5251936131500646E8, -3.746685949111144E7, 1.5160668038148597E8, -4.852577939892976E8, 1.7443933581116918E8, 1.0531243116544546E8, -1.1895546897881567E9, 6.341783697427113E7, -1.0726206418402787E8, -1.0839573118002522E9,
				}),
			},
			lineCentroid: geom.Coord{-3.746685949111144E7, 1.5160668038148597E8, -4.852577939892976E8},
		},
		{
			desc: "Randomly Generated 4",
			lines: []*geom.LineString{
				geom.NewLineStringFlat(geom.XYZ, []float64{
					-3773352.7330737715, -5.342411661592115E7, 5.3125486291572404E8, -3.158050487041971E8, 1.617940676945237E8, 5.657757480926482E7, -3773352.7330737715, -5.342411661592115E7, 5.3125486291572404E8,
				}),
			},
			lineCentroid: geom.Coord{-3.158050487041971E8, 1.617940676945237E8, 5.657757480926482E7},
		},
		{
			desc: "Randomly Generated 5",
			lines: []*geom.LineString{
				geom.NewLineStringFlat(geom.XYZ, []float64{
					-6.955916118828537E8, -7.54771571550343E8, -1.0941742587086952E8, -1.1781096802518074E8, -1.7082976090687165E8, -1.201009454243481E9, 6.354062301153038E8, 2.0204896090364275E9, 1.962046862159836E9, -2.8653780223221374E8, -9.464234534114834E8, -9.911874544497331E8, -6.955916118828537E8, -7.54771571550343E8, -1.0941742587086952E8,
				}),
			},
			lineCentroid: geom.Coord{-1.1781096802518074E8, -1.7082976090687165E8, -1.201009454243481E9},
		},
	} {
		verifyInteriorPointBasicLines(t, tc)
		verifyInteriorPointMultiLine(t, tc)
		verifyInteriorPointLinearRing(t, tc)
	}
}

func verifyInteriorPointBasicLines(t *testing.T, tc lineDataType) {
	interiorPoint := xy.LinesInteriorPoint(tc.lines[0], tc.lines[1:]...)

	if !reflect.DeepEqual(interiorPoint, tc.lineCentroid) {
		t.Errorf("Test %v Failed.  Expected \n\t%v but was\n\t%v", tc.desc, tc.lineCentroid, interiorPoint)
	}
}

func verifyInteriorPointMultiLine(t *testing.T, tc lineDataType) {
	ends := []int{}
	coords := []float64{}
	for _, line := range tc.lines {
		coords = append(coords, line.FlatCoords()...)
		ends = append(ends, len(coords))
	}
	multiline := geom.NewMultiLineStringFlat(tc.lines[0].Layout(), coords, ends)
	interiorPoint := xy.MultiLineInteriorPoint(multiline)

	if !reflect.DeepEqual(interiorPoint, tc.lineCentroid) {
		t.Errorf("Test %v (Multiline) Failed.  Expected \n\t%v but was\n\t%v", tc.desc, tc.lineCentroid, interiorPoint)
	}
}

func verifyInteriorPointLinearRing(t *testing.T, tc lineDataType) {
	rings := []*geom.LinearRing{}

	for _, line := range tc.lines {
		coords := make([]float64, len(line.FlatCoords()))
		copy(coords, line.FlatCoords())
		rings = append(rings, geom.NewLinearRingFlat(line.Layout(), coords))
	}
	interiorPoint := xy.LinearRingsInteriorPoint(rings[0], rings[1:]...)

	if !reflect.DeepEqual(interiorPoint, tc.lineCentroid) {
		t.Errorf("Test %v (LinearRing) Failed.  Expected \n\t%v but was\n\t%v", tc.desc, tc.lineCentroid, interiorPoint)
	}
}
