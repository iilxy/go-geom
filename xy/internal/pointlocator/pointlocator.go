package pointlocator

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/boundary"
	"github.com/twpayne/go-geom/xy/internal"
	"github.com/twpayne/go-geom/xy/internal/lineintersector"
	"github.com/twpayne/go-geom/xy/internal/raycrossing"
	"github.com/twpayne/go-geom/xy/location"
)

// LocatePointOnGeom computes the topological relationship ({@link Location}) of a single point
// to a Geometry.
// It handles both single-element and multi-element Geometries.
// The algorithm for multi-part Geometries takes into account the SFS Boundary Determination Rule.
func LocatePointOnGeom(boundaryRule boundary.NodeRule, point geom.Coord, geometry geom.T) location.Type {
	locator := pointLocator{
		boundaryRule:  boundaryRule,
		isIn:          false,
		numBoundaries: 0,
	}

	return locator.locatePointOnGeom(point, geometry)
}

type pointLocator struct {
	boundaryRule boundary.NodeRule
	// true if the point lies in or on any Geometry element
	isIn bool
	// the number of sub-elements whose boundaries the point lies in
	numBoundaries int
}

func (loc *pointLocator) locatePointOnGeom(point geom.Coord, geometry geom.T) location.Type {
	if len(geometry.FlatCoords()) == 0 {
		return location.Exterior
	}

	switch g := geometry.(type) {
	case *geom.LineString:
		return loc.locatePointOnLine(point, g)
	case *geom.Polygon:
		return loc.locatePointOnPolygon(point, g)
	}

	loc.computeLocation(point, geometry)
	if loc.boundaryRule.IsInBoundary(loc.numBoundaries) {
		return location.Boundary
	} else if loc.numBoundaries > 0 || loc.isIn {
		return location.Interior
	}

	return location.Exterior
}

func (loc *pointLocator) computeLocation(point geom.Coord, geometry geom.T) {

	switch g := geometry.(type) {
	case *geom.Point:
		loc.updateLocationInfo(loc.locatePointOnPoint(point, g))
	case *geom.LineString:
		loc.updateLocationInfo(loc.locatePointOnLine(point, g))
	case *geom.MultiLineString:
		for i := 0; i < g.NumLineStrings(); i++ {
			l := g.LineString(i)
			loc.updateLocationInfo(loc.locatePointOnLine(point, l))
		}
	case *geom.Polygon:
		loc.updateLocationInfo(loc.locatePointOnPolygon(point, g))
	case *geom.MultiPolygon:
		for i := 0; i < g.NumPolygons(); i++ {
			poly := g.Polygon(i)
			loc.updateLocationInfo(loc.locatePointOnPolygon(point, poly))
		}
	}
}

func (loc *pointLocator) updateLocationInfo(currentLoc location.Type) {
	if currentLoc == location.Interior {
		loc.isIn = true
	}
	if currentLoc == location.Boundary {
		loc.numBoundaries++
	}
}

func (loc *pointLocator) locatePointOnPoint(point geom.Coord, line *geom.Point) location.Type {
	// no point in doing envelope test, since equality test is just as fast
	if internal.Equal(point, 0, line.FlatCoords(), 0) {
		return location.Interior
	}
	return location.Exterior
}

func (loc *pointLocator) locatePointOnLine(point geom.Coord, line *geom.LineString) location.Type {
	// bounding-box check
	coords := line.FlatCoords()
	stride := line.Stride()
	if !line.Bounds().OverlapsPoint(geom.XY, point) {
		return location.Exterior
	}

	lineClosed := internal.Equal(coords, 0, coords, len(coords)-stride)
	if !lineClosed {
		if internal.Equal(point, 0, coords, 0) || internal.Equal(point, 0, coords, len(coords)-stride) {
			return location.Boundary
		}
	}

	if lineintersector.IsOnLine(line.Layout(), point, coords) {
		return location.Interior
	}
	return location.Exterior
}

func (loc *pointLocator) locatePointOnPolygon(point geom.Coord, poly *geom.Polygon) location.Type {
	coords := poly.FlatCoords()
	if len(coords) == 0 {
		return location.Exterior
	}

	shell := poly.LinearRing(0)

	shellLoc := loc.locateInPolygonRing(point, shell)
	if shellLoc == location.Exterior {
		return location.Exterior
	}
	if shellLoc == location.Boundary {
		return location.Boundary
	}

	// now test if the point lies in or on the holes
	for i := 1; i < poly.NumLinearRings(); i++ {
		hole := poly.LinearRing(i)
		holeLoc := loc.locateInPolygonRing(point, hole)

		if holeLoc == location.Interior {
			return location.Exterior
		}
		if holeLoc == location.Boundary {
			return location.Boundary
		}
	}

	return location.Interior
}

func (loc *pointLocator) locateInPolygonRing(p geom.Coord, ring *geom.LinearRing) location.Type {
	// bounding-box check
	if !ring.Bounds().OverlapsPoint(geom.XY, p) {
		return location.Exterior
	}

	return raycrossing.LocatePointInRing(ring.Layout(), p, ring.FlatCoords())
}
