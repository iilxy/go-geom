package geom_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/twpayne/go-geom"
)

func TestBuilder_Point(t *testing.T) {
	point, err := geom.Build(geom.XY).StartPoint(1, 1).Point()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY.StartPoint(1, 1).Point() should not have given an error but got: %v", err)
	}

	if point == nil {
		t.Fatal("geom.Build(geom.XY.StartPoint(1, 1).Point() should have returned a non-nil point")
	}

	if point.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY.StartPoint(1, 1).Point() has layout %v wanted %v", point.Layout(), geom.XY)
	}

	if !reflect.DeepEqual(point.FlatCoords(), []float64{1, 1}) {
		t.Errorf("geom.Build(geom.XY.StartPoint(1, 1).Point() coords: %v wanted %v", point.FlatCoords(), []float64{1, 1})
	}
}

func TestBuilder_Point_InsufficientOrdinals(t *testing.T) {
	point, err := geom.Build(geom.XYZ).StartPoint(1, 1).Point()

	if err == nil {
		t.Error("geom.Build(geom.XYZ).StartPoint(1, 1).Point() should have returned an error")
	}

	if !strings.Contains(err.Msg, "correct number of ordinates") {
		t.Errorf("geom.Build(geom.XYZ).StartPoint(1, 1).Point() returned the wrong error message: %v", err)
	}

	if point != nil {
		t.Errorf("geom.Build(geom.XYZ).StartPoint(1, 1).Point() should not have returned a point. Got: %v", point)
	}
}

func TestBuilder_Point_MissingGoTo(t *testing.T) {
	point, err := geom.Build(geom.XYZ).Point()

	if err == nil {
		t.Error("geom.Build(geom.XYZ).Point() should have returned an error")
	}

	if !strings.Contains(err.Msg, "A Point cannot be created because") {
		t.Errorf("geom.Build(geom.XYZ).Point() returned the wrong error message: %v", err)
	}

	if point != nil {
		t.Errorf("geom.Build(geom.XYZ).Point() should not have returned a point. Got: %v", point)
	}
}

func TestBuilder_Point_LastPoint(t *testing.T) {
	point, err := geom.Build(geom.XY).StartPoint(1, 2).StartPoint(3, 4).Point()

	if err != nil {
		t.Errorf("geom.Build(geom.XY).StartPoint(1, 2).StartPoint(3, 4).Point() should not have returned an error. Got: %v", err)
	}

	if point == nil {
		t.Error("geom.Build(geom.XY).StartPoint(1, 2).StartPoint(3, 4).Point() should have returned a point.")
	}

	expectedCoords := []float64{3, 4}
	if !reflect.DeepEqual(point.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartPoint(1, 2).StartPoint(3, 4).Point() has coords: %v wanted: %v", point.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_MultiPoint(t *testing.T) {
	mp, err := geom.Build(geom.XY).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(true)

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint() should not have given an error but got: %v", err)
	}

	if mp == nil {
		t.Fatal("geom.Build(geom.XY).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint() should have returned a non-nil point")
	}

	if mp.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint() has layout %v wanted %v", mp.Layout(), geom.XY)
	}

	wanted := []float64{1, 1, 3, 3, 2, 2}
	if !reflect.DeepEqual(mp.FlatCoords(), wanted) {
		t.Errorf("geom.Build(geom.XY).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint() coords: %v wanted %v", mp.FlatCoords(), wanted)
	}
}

func TestBuilder_MultiPoint_IgnoreLine(t *testing.T) {
	mp, err := geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(false)

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(false) should not have given an error but got: %v", err)
	}

	if mp == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(false) should have returned a non-nil geometry")
	}

	if mp.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(false) has layout %v wanted %v", mp.Layout(), geom.XY)
	}

	wanted := []float64{1, 1, 3, 3, 2, 2}
	if !reflect.DeepEqual(mp.FlatCoords(), wanted) {
		t.Errorf("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(false) coords: %v wanted %v", mp.FlatCoords(), wanted)
	}
}

func TestBuilder_MultiPoint_NotAllPoints(t *testing.T) {
	mp, err := geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(true)

	if err == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(true) should have given an error")
	}

	if !strings.Contains(err.Msg, "was not a point") {
		t.Errorf("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(true) returned the wrong error message: %v", err)
	}

	if mp != nil {
		t.Fatal("geom.Build(geom.XY).StartLine(-1, -1).LineTo(-2, -2).StartPoint(1, 1).StartPoint(3, 3).StartPoint(2, 2).MultiPoint(true) should not have geometry")
	}
}

func TestBuilder_Line(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{1, 1, 10, 10}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_Line_Continue(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine(20, 20).LineTo(30, 0).LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine(20, 20).LineTo(30, 0).LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine(20, 20).LineTo(30, 0).LineString() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).GoTo(1, 1).LineTo(10, 10).GoTo(20,20).LineTo(30,0)LineString() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{20, 20, 30, 0}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).GoTo(1, 1).LineTo(10, 10).GoTo(20,20).LineTo(30,0)LineString() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_Line_Continue2(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine().LineTo(30, 0).LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine().LineTo(30, 0).LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine().LineTo(30, 0).LineString() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine().LineTo(30, 0).LineString() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{10, 10, 30, 0}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine().LineTo(30, 0).LineString() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_Line_WrongLayout(t *testing.T) {
	line, err := geom.Build(geom.XYZ).StartLine(1, 1).LineTo(10, 10).LineString()

	if err == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() should have produced an error.")
	}

	if !strings.Contains(err.Msg, "correct number of ordinates") {
		t.Errorf("geom.Build(geom.XYZ).StartLine(1, 1).LineTo(10, 10).LineString() returned the wrong error message: %v", err)
	}

	if line != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() should not have produced line. Got: %v", line)
	}
}

func TestBuilder_Line_NoCoordsOnStart(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine().LineTo(10, 10).LineString()

	if err == nil {
		t.Fatal("geom.Build(geom.XYZ).StartLine().LineTo(10, 10).LineString() should have produced an error.")
	}

	if !strings.Contains(err.Msg, "geometry has already been created") {
		t.Errorf("geom.Build(geom.XY).StartLine().LineTo(10, 10).LineString() returned the wrong error message: %v", err)
	}

	if line != nil {
		t.Fatalf("geom.Build(geom.XYZ).StartLine().LineTo(10, 10).LineString() should not have produced line. Got: %v", line)
	}
}

func TestBuilder_Close(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{1, 1, 10, 10, 10, 20, 1, 1}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_Close_OnlyCloseWhenRequired(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).LineTo(1, 1).CloseRing().LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{1, 1, 10, 10, 10, 20, 1, 1}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 20).CloseRing().LineString() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_Close_InsufficientCoords(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).CloseRing().LineString()

	if err == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).CloseRing().LineString() should have produced an error")
	}

	if !strings.Contains(err.Msg, "at least 3 coordinates") {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).CloseRing().LineString() returned the wrong error message: %v", err)
	}

	if line != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).CloseRing().LineString() should not have produced line: %v", line)
	}
}

func TestBuilder_Close_NeedsStart(t *testing.T) {
	line, err := geom.Build(geom.XY).CloseRing().LineString()

	if err == nil {
		t.Fatal("geom.Build(geom.XY).CloseRing().LineString() should have produced an error")
	}

	if !strings.Contains(err.Msg, "must be started") {
		t.Errorf("geom.Build(geom.XY).CloseRing().LineString() returned the wrong error message: %v", err)
	}

	if line != nil {
		t.Fatalf("geom.Build(geom.XY).CloseRing().LineString() should not have produced line: %v", line)
	}
}

func TestBuilder_LinearRing(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 15).LinearRing()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 15).LinearRing() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 15).LinearRing() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 15).LinearRing() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{1, 1, 10, 10, 10, 15, 1, 1}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineTo(10, 15).LinearRing() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_AddLineSegments(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).AddLineSegments(geom.Coord{10, 10}, geom.Coord{10, 20}).LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).AddLineSegments(geom.Coord{10, 10}, geom.Coord{10, 20}).LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatal("geom.Build(geom.XY).StartLine(1, 1).AddLineSegments(geom.Coord{10, 10}, geom.Coord{10, 20}).LineString() should have produced line but did not")
	}

	if line.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).AddLineSegments(geom.Coord{10, 10}, geom.Coord{10, 20}).LineString() produced a line with the wrong layout.  Was %v wanted %v", line.Layout(), geom.XY)
	}

	expectedCoords := []float64{1, 1, 10, 10, 10, 20}
	if !reflect.DeepEqual(line.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).StartLine(1, 1).AddLineSegments(geom.Coord{10, 10}, geom.Coord{10, 20}).LineString() produced a line with the wrong coords.  Was %v wanted %v", line.FlatCoords(), expectedCoords)
	}
}

func TestBuilder_MultiLineString(t *testing.T) {
	ml, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 0).LineTo(10, 10).
		StartLine(2, 2).LineTo(2, 4).LineTo(4, 4).
		MultiLineString(false)

	if err != nil {
		t.Fatalf("geom.Build(geom.XY)....MultiLineString() should not have produced an error: %v", err)
	}

	if ml == nil {
		t.Fatal("geom.Build(geom.XY).....MultiLineString() should have produced MultiLineString but did not")
	}

	if ml.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).....MultiLineString() produced a MultiLineString with the wrong layout.  Was %v wanted %v", ml.Layout(), geom.XY)
	}

	expectedCoords := []float64{
		1, 1, 10, 0, 10, 10,
		2, 2, 2, 4, 4, 4}

	if !reflect.DeepEqual(ml.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).....MultiLineString() produced a MultiLineString with the wrong coords.  Was %v wanted %v", ml.FlatCoords(), expectedCoords)
	}
	expectedEnds := []int{6, 12}
	if !reflect.DeepEqual(ml.Ends(), expectedEnds) {
		t.Errorf("geom.Build(geom.XY).....MultiLineString() produced a MultiLineString with the wrong ends.  Was %v wanted %v", ml.Ends(), expectedEnds)
	}
}

func TestBuilder_Polygon(t *testing.T) {
	polygon, err := geom.Build(geom.XY).StartPolygon(1, 1).LineTo(10, 0).LineTo(10, 10).
		StartHole(2, 2).LineTo(2, 4).LineTo(4, 4).
		StartHole(6, 6).LineTo(6, 8).LineTo(8, 8).
		Polygon()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY)....Polygon() should not have produced an error: %v", err)
	}

	if polygon == nil {
		t.Fatal("geom.Build(geom.XY).....Polygon() should have produced polygon but did not")
	}

	if polygon.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).....Polygon() produced a polygon with the wrong layout.  Was %v wanted %v", polygon.Layout(), geom.XY)
	}

	expectedCoords := []float64{
		1, 1, 10, 0, 10, 10, 1, 1,
		2, 2, 2, 4, 4, 4, 2, 2,
		6, 6, 6, 8, 8, 8, 6, 6}

	if !reflect.DeepEqual(polygon.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).....Polygon() produced a polygon with the wrong coords.  Was %v wanted %v", polygon.FlatCoords(), expectedCoords)
	}
	expectedEnds := []int{8, 16, 24}
	if !reflect.DeepEqual(polygon.Ends(), expectedEnds) {
		t.Errorf("geom.Build(geom.XY).....Polygon() produced a polygon with the wrong ends.  Was %v wanted %v", polygon.Ends(), expectedEnds)
	}
}

func TestBuilder_MultiPolygon_AllGeomsFalse(t *testing.T) {
	multipolygon, err := geom.Build(geom.XY).
		StartPoint(1, 1).
		StartPolygon(1, 1).LineTo(10, 0).LineTo(10, 10).
		StartPolygon(-1, -1).LineTo(-10, 0).LineTo(-10, -10).
		StartPolygon(21, 21).LineTo(30, 20).LineTo(30, 30).
		MultiPolygon(false)

	validationsForMultiPolygon(multipolygon, err, t)
}

func TestBuilder_MultiPolygon_AllGeomsTrue(t *testing.T) {
	multipolygon, err := geom.Build(geom.XY).
		StartPolygon(1, 1).LineTo(10, 0).LineTo(10, 10).
		StartPolygon(-1, -1).LineTo(-10, 0).LineTo(-10, -10).
		StartPolygon(21, 21).LineTo(30, 20).LineTo(30, 30).
		MultiPolygon(true)

	validationsForMultiPolygon(multipolygon, err, t)
}

func validationsForMultiPolygon(multipolygon *geom.MultiPolygon, err *geom.BuilderError, t *testing.T) {
	if err != nil {
		t.Fatalf("geom.Build(geom.XY)....MultiPolygon() should not have produced an error: %v", err)
	}

	if multipolygon == nil {
		t.Fatal("geom.Build(geom.XY).....MultiPolygon() should have produced MultiPolygon but did not")
	}

	if multipolygon.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY).....MultiPolygon() produced a MultiPolygon with the wrong layout.  Was %v wanted %v", multipolygon.Layout(), geom.XY)
	}
	expectedCoords := []float64{
		1, 1, 10, 0, 10, 10, 1, 1,
		-1, -1, -10, 0, -10, -10, -1, -1,
		21, 21, 30, 20, 30, 30, 21, 21,
	}

	if !reflect.DeepEqual(multipolygon.FlatCoords(), expectedCoords) {
		t.Errorf("geom.Build(geom.XY).....MultiPolygon() produced a MultiPolygon with the wrong coords.  Was %v wanted %v", multipolygon.FlatCoords(), expectedCoords)
	}

	if multipolygon.Ends() != nil {
		t.Errorf("geom.Build(geom.XY).....MultiPolygon() produced a MultiPolygon with the wrong ends.  Was %v wanted %v", multipolygon.Ends(), nil)
	}

	expectedEndss := [][]int{[]int{8}, []int{16}, []int{24}}
	if !reflect.DeepEqual(multipolygon.Endss(), expectedEndss) {
		t.Errorf("geom.Build(geom.XY).....MultiPolygon() produced a MultiPolygon with the wrong endss.  Was %v wanted %v", multipolygon.Endss(), expectedEndss)
	}
}

func TestBuilder_MultiPolygon_Error_IncorrectMemberType(t *testing.T) {
	multipolygon, err := geom.Build(geom.XY).
		StartPoint(1, 1).
		StartPolygon(1, 1).LineTo(10, 0).LineTo(10, 10).
		StartPolygon(-1, -1).LineTo(-10, 0).LineTo(-10, -10).
		StartPolygon(21, 21).LineTo(30, 20).LineTo(30, 30).
		MultiPolygon(true)

	if err == nil {
		t.Fatal("geom.Build(geom.XY)....MultiPolygon() should have produced an error but it did not")
	}

	if multipolygon != nil {
		t.Fatalf("geom.Build(geom.XY).....MultiPolygon() should not have produced MultiPolygon: %v", multipolygon)
	}

}
