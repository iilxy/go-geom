package pointlocator_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/boundary"
	"github.com/twpayne/go-geom/xy/internal/pointlocator"
	"github.com/twpayne/go-geom/xy/location"
	"testing"
)

type locatePointTestData struct {
	desc     string
	point    geom.Coord
	geometry geom.T
	nodeRule boundary.NodeRule
	result   location.Type
}

func TestLocatePointOnGeomPoints(t *testing.T) {
	for _, tc := range []locatePointTestData{
		{
			desc:     "point beside point",
			point:    geom.Coord{0, 0},
			geometry: geom.NewPointFlat(geom.XY, []float64{1, 1}),
			result:   location.Exterior,
		},
		{
			desc:     "point on point",
			point:    geom.Coord{0, 0},
			geometry: geom.NewPointFlat(geom.XY, []float64{0, 0}),
			result:   location.Interior,
		},
	} {
		onGeom := pointlocator.LocatePointOnGeom(nodeRule(tc), tc.point, tc.geometry)

		if onGeom != tc.result {
			t.Errorf("Test '%s' failed.  Expected %v but was %v", tc.desc, tc.result, onGeom)
		}
	}
}

func TestLocatePointOnGeomLineString(t *testing.T) {
	for _, tc := range []locatePointTestData{
		{
			desc:     "point beside line",
			point:    geom.Coord{0, 0},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{1, 0, 1, 1}),
			result:   location.Exterior,
		},
		{
			desc:     "point on line",
			point:    geom.Coord{0, 0},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{-1, -1, 1, 1}),
			result:   location.Interior,
		},
		{
			desc:     "point on closed line OGC nodeRule",
			point:    geom.Coord{0, 1},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}),
			result:   location.Interior,
		},
		{
			desc:     "point on closed line End Point Boundary nodeRule",
			point:    geom.Coord{0, 1},
			nodeRule: boundary.EndPointBoundaryNodeRule{},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}),
			result:   location.Interior,
		},
		{
			desc:     "point on line endpoint",
			point:    geom.Coord{0, 0},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 1}),
			result:   location.Boundary,
		},
	} {
		onGeom := pointlocator.LocatePointOnGeom(nodeRule(tc), tc.point, tc.geometry)

		if onGeom != tc.result {
			t.Errorf("Test '%s' failed.  Expected %v but was %v", tc.desc, tc.result, onGeom)
		}
	}
}

func TestLocatePointOnGeomLinearRing(t *testing.T) {
	for _, tc := range []locatePointTestData{
		{
			desc:     "point beside linear ring",
			point:    geom.Coord{0, 0},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}),
			result:   location.Exterior,
		},
		{
			desc:     "point in linear ring",
			point:    geom.Coord{1.5, 1.5},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}),
			result:   location.Exterior,
		},
		{
			desc:     "point on border of linear ring OGC node rule",
			point:    geom.Coord{1, 1},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}),
			result:   location.Interior,
		},
		{
			desc:     "point on border of linear ring EndPoint Boundary NodeRule",
			point:    geom.Coord{1, 1},
			nodeRule: boundary.EndPointBoundaryNodeRule{},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}),
			result:   location.Interior,
		},
	} {
		onGeom := pointlocator.LocatePointOnGeom(nodeRule(tc), tc.point, tc.geometry)

		if onGeom != tc.result {
			t.Errorf("Test '%s' failed.  Expected %v but was %v", tc.desc, tc.result, onGeom)
		}
	}
}

func TestLocatePointOnGeomMultiLineString(t *testing.T) {
	for _, tc := range []locatePointTestData{
		{
			desc:     "point between lines",
			point:    geom.Coord{0, 0},
			geometry: geom.NewMultiLineStringFlat(geom.XY, []float64{1, 1, 1, -1, -1, 1, -1, -1}, []int{4, 8}),
			result:   location.Exterior,
		},
		{
			desc:     "point on one line",
			point:    geom.Coord{-1, 0},
			geometry: geom.NewMultiLineStringFlat(geom.XY, []float64{1, 1, 1, -1, -1, 1, -1, -1}, []int{4, 8}),
			result:   location.Interior,
		},
		{
			desc:     "point on line endpoint",
			point:    geom.Coord{-1, 1},
			geometry: geom.NewMultiLineStringFlat(geom.XY, []float64{1, 1, 1, -1, -1, 1, -1, -1}, []int{4, 8}),
			result:   location.Boundary,
		},
	} {
		onGeom := pointlocator.LocatePointOnGeom(nodeRule(tc), tc.point, tc.geometry)

		if onGeom != tc.result {
			t.Errorf("Test '%s' failed.  Expected %v but was %v", tc.desc, tc.result, onGeom)
		}
	}
}

func TestLocatePointOnGeomPolygon(t *testing.T) {
	for _, tc := range []locatePointTestData{
		{
			desc:     "point beside polygon",
			point:    geom.Coord{0, 0},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}, []int{10}),
			result:   location.Exterior,
		},
		{
			desc:     "point in polygon",
			point:    geom.Coord{1.5, 1.5},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}, []int{10}),
			result:   location.Interior,
		},
		{
			desc:     "point on polygon boundary OGC nodeRule",
			point:    geom.Coord{1, 1},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}, []int{10}),
			result:   location.Boundary,
		},
		{
			desc:     "point on polygon boundary EndPoint Boundary NodeRule",
			point:    geom.Coord{1, 1},
			nodeRule: boundary.EndPointBoundaryNodeRule{},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}, []int{10}),
			result:   location.Boundary,
		},
		{
			desc:  "point in polygon hole",
			point: geom.Coord{1.5, 1.5},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				1.25, 1.25, 1.75, 1.25, 1.75, 1.75, 1.25, 1.75, 1.25, 1.25}, []int{10, 20}),
			result: location.Exterior,
		},
		{
			desc:     "point in polygon beside hole",
			point:    geom.Coord{1.1, 1.1},
			nodeRule: boundary.EndPointBoundaryNodeRule{},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				1.25, 1.25, 1.75, 1.25, 1.75, 1.75, 1.25, 1.75, 1.25, 1.25}, []int{10, 20}),
			result: location.Interior,
		},
	} {
		onGeom := pointlocator.LocatePointOnGeom(nodeRule(tc), tc.point, tc.geometry)

		if onGeom != tc.result {
			t.Errorf("Test '%s' failed.  Expected %v but was %v", tc.desc, tc.result, onGeom)
		}
	}
}

func TestLocatePointOnGeomMultiPolygon(t *testing.T) {
	for _, tc := range []locatePointTestData{
		{
			desc:  "point between multi-polygon",
			point: geom.Coord{0, 0},
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				-1, -1, -2, -1, -2, -2, -1, -2, -1, -1,
			}, [][]int{
				[]int{10},
				[]int{20}}),
			result: location.Exterior,
		},
		{
			desc:  "point in multi-polygon",
			point: geom.Coord{1.5, 1.5},
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				-1, -1, -2, -1, -2, -2, -1, -2, -1, -1,
			}, [][]int{
				[]int{10},
				[]int{20}}),
			result: location.Interior,
		},
		{
			desc:  "point in hole multi-polygon",
			point: geom.Coord{1.5, 1.5},
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				1.25, 1.25, 1.75, 1.25, 1.75, 1.75, 1.25, 1.75, 1.25, 1.25,
				-1, -1, -2, -1, -2, -2, -1, -2, -1, -1,
			}, [][]int{
				[]int{10, 20},
				[]int{20}}),
			result: location.Exterior,
		},
		{
			desc:  "point in overlapping multi-polygon",
			point: geom.Coord{1.5, 1.5},
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				0, 0, 1.75, 0, 1.75, 1.75, 0, 1.75, 0, 0,
			}, [][]int{
				[]int{10},
				[]int{20}}),
			result: location.Interior,
		},
		{
			desc:  "point in hole of overlapping multi-polygon, hole is overlapped",
			point: geom.Coord{1.5, 1.5},
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{
				1, 1, 2, 1, 2, 2, 1, 2, 1, 1,
				1.25, 1.25, 1.75, 1.25, 1.75, 1.75, 1.25, 1.75, 1.25, 1.25,
				0, 0, 1.75, 0, 1.75, 1.75, 0, 1.75, 0, 0,
			}, [][]int{
				[]int{10},
				[]int{20}}),
			result: location.Interior,
		},
	} {
		onGeom := pointlocator.LocatePointOnGeom(nodeRule(tc), tc.point, tc.geometry)

		if onGeom != tc.result {
			t.Errorf("Test '%s' failed.  Expected %v but was %v", tc.desc, tc.result, onGeom)
		}
	}
}

func nodeRule(tc locatePointTestData) boundary.NodeRule {
	if tc.nodeRule != nil {
		return tc.nodeRule
	}
	return boundary.Mod2BoundaryNodeRule{}
}
