package geom_test

import (
	"testing"

	"reflect"

	"fmt"

	"github.com/twpayne/go-geom"
)

func TestBuilder_Point(t *testing.T) {
	point, err := geom.Build(geom.XY).GoTo(1, 1).Point()

	if err != nil {
		t.Fatalf("geom.Build(geom.XY.GoTo(1, 1).Point() should not have given an error but got: %v", err)
	}

	if point == nil {
		t.Fatal("geom.Build(geom.XY.GoTo(1, 1).Point() should have returned a non-nil point")
	}

	if point.Layout() != geom.XY {
		t.Errorf("geom.Build(geom.XY.GoTo(1, 1).Point() has layout %v wanted %v", point.Layout(), geom.XY)
	}

	if !reflect.DeepEqual(point.FlatCoords(), []float64{1, 1}) {
		t.Errorf("geom.Build(geom.XY.GoTo(1, 1).Point() coords: %v wanted %v", point.FlatCoords(), []float64{1, 1})
	}
}

func TestBuilder_Point_InsufficientOrdinals(t *testing.T) {
	point, err := geom.Build(geom.XYZ).GoTo(1, 1).Point()

	if err == nil {
		t.Error("geom.Build(geom.XYZ).GoTo(1, 1).Point() should have returned an error")
	}

	if point != nil {
		t.Errorf("geom.Build(geom.XYZ).GoTo(1, 1).Point() should not have returned a point. Got: %v", point)
	}
}

func TestBuilder_Point_MissingGoTo(t *testing.T) {
	point, err := geom.Build(geom.XYZ).Point()

	if err == nil {
		t.Error("geom.Build(geom.XYZ).Point() should have returned an error")
	}

	if point != nil {
		t.Errorf("geom.Build(geom.XYZ).Point() should not have returned a point. Got: %v", point)
	}
}

func TestBuilder_Point_NotAPoint(t *testing.T) {
	point, err := geom.Build(geom.XY).GoTo(1, 2).GoTo(3, 4).Point()

	fmt.Println(err)
	if err == nil {
		t.Error("geom.Build(geom.XY).GoTo(1, 2).GoTo(3, 4).Point() should have returned an error")
	}

	if point != nil {
		t.Errorf("geom.Build(geom.XY).GoTo(1, 2).GoTo(3, 4).Point() should not have returned a point. Got: %v", point)
	}
}
