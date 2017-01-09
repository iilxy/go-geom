package geom_test

import (
	"testing"

	"reflect"

	"strings"

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

	if !strings.Contains(err.Msg, "A point cannot be created") {
		t.Errorf("geom.Build(geom.XYZ).StartPoint(1, 1).Point() returned the wrong error message: %v", err)
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

func TestBuilder_Line(t *testing.T) {
	line, err := geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() should not have produced an error: %v", err)
	}

	if line == nil {
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).LineString() should have produced line but did not")
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
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine(20, 20).LineTo(30, 0).LineString() should have produced line but did not")
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
		t.Fatalf("geom.Build(geom.XY).StartLine(1, 1).LineTo(10, 10).StartLine().LineTo(30, 0).LineString() should have produced line but did not")
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
		t.Errorf("geom.Build(geom.XYZ).StartPoint(1, 1).Point() returned the wrong error message: %v", err)
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
		t.Errorf("geom.Build(geom.XYZ).StartPoint(1, 1).Point() returned the wrong error message: %v", err)
	}

	if line != nil {
		t.Fatalf("geom.Build(geom.XYZ).StartLine().LineTo(10, 10).LineString() should not have produced line. Got: %v", line)
	}
}
