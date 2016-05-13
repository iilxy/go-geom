package lineintersector_test

import (
	"testing"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/internal/lineintersector"
)

func TestIsOnLinePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("This test is supposed to panic")
		}
		// good panic was expected
	}()

	lineintersector.IsOnLine(geom.XY, geom.Coord{0, 0}, []float64{0, 0})
}

func TestIsOnLine(t *testing.T) {
	for i, tc := range []struct {
		desc         string
		p            geom.Coord
		lineSegments []float64
		layout       geom.Layout
		intersects   bool
	}{
		{
			desc:         "Point on center of line",
			p:            geom.Coord{0, 0},
			lineSegments: []float64{-1, 0, 1, 0},
			layout:       geom.XY,
			intersects:   true,
		},
		{
			desc:         "Point not on line",
			p:            geom.Coord{0, 0},
			lineSegments: []float64{-1, 1, 1, 0},
			layout:       geom.XY,
			intersects:   false,
		},
		{
			desc:         "Point not on second line segment",
			p:            geom.Coord{0, 0},
			lineSegments: []float64{-1, 1, 1, 0, -1, 0},
			layout:       geom.XY,
			intersects:   true,
		},
		{
			desc:         "Point not on any line segments",
			p:            geom.Coord{0, 0},
			lineSegments: []float64{-1, 1, 1, 0, 2, 0},
			layout:       geom.XY,
			intersects:   false,
		},
		{
			desc:         "Point in unclosed ring",
			p:            geom.Coord{0, 0},
			lineSegments: []float64{-1, 1, 1, 1, 1, -1, -1, -1, -1, 1.00000000000000000000000000001},
			layout:       geom.XY,
			intersects:   false,
		},
		{
			desc:         "Point in ring",
			p:            geom.Coord{0, 0},
			lineSegments: []float64{-1, 1, 1, 1, 1, -1, -1, -1, -1, 1},
			layout:       geom.XY,
			intersects:   false,
		},
	} {
		if tc.intersects != lineintersector.IsOnLine(tc.layout, tc.p, tc.lineSegments) {
			t.Errorf("Test '%v' (%v) failed: expected \n%v but was \n%v", i+1, tc.desc, tc.intersects, !tc.intersects)
		}
	}
}
