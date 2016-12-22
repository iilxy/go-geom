// Package geom implements efficient geometry types for geospatial
// applications.
package geom

import (
	"errors"
	"fmt"
	"github.com/twpayne/go-geom/ogc"
	"math"
)

// A Layout describes the meaning of an N-dimensional coordinate. Layout(N) for
// N > 4 is a valid layout, in which case the first dimensions are interpreted
// to be X, Y, Z, and M and extra dimensions have no special meaning.  M values
// are considered part of a linear referencing system (e.g. classical time or
// distance along a path). 1-dimensional layouts are not supported.
type Layout int

const (
	// NoLayout is an unknown layout
	NoLayout Layout = iota
	// XY is a 2D layout (X and Y)
	XY
	// XYZ is 3D layout (X, Y, and Z)
	XYZ
	// XYM is a 2D layout with an M value
	XYM
	// XYZM is a 3D layout with an M value
	XYZM
)

// An ErrLayoutMismatch is returned when geometries with different layouts
// cannot be combined.
type ErrLayoutMismatch struct {
	Got  Layout
	Want Layout
}

func (e ErrLayoutMismatch) Error() string {
	return fmt.Sprintf("geom: layout mismatch, got %s, want %s", e.Got, e.Want)
}

// An ErrStrideMismatch is returned when the stride does not match the expected
// stride.
type ErrStrideMismatch struct {
	Got  int
	Want int
}

func (e ErrStrideMismatch) Error() string {
	return fmt.Sprintf("geom: stride mismatch, got %d, want %d", e.Got, e.Want)
}

// An ErrUnsupportedLayout is returned when the requested layout is not
// supported.
type ErrUnsupportedLayout Layout

func (e ErrUnsupportedLayout) Error() string {
	return fmt.Sprintf("geom: unsupported layout %s", Layout(e))
}

// An ErrUnsupportedType is returned when the requested type is not supported.
type ErrUnsupportedType struct {
	Value interface{}
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("geom: unsupported type %T", e.Value)
}

// A Coord represents an N-dimensional coordinate.
type Coord []float64

// Clone returns a deep copy of c.
func (c Coord) Clone() Coord {
	clone := make(Coord, len(c))
	copy(clone, c)
	return clone
}

// X returns the x coordinate of the coordinate. X is assumed to be the first
// ordinate.
func (c Coord) X() float64 {
	return c[0]
}

// Y returns the x coordinate of the coordinate. Y is assumed to be the second
// ordinate.
func (c Coord) Y() float64 {
	return c[1]
}

// Set copies the ordinate data from the other coord to this coord.
func (c Coord) Set(other Coord) {
	copy(c, other)
}

// Equal compares that all ordinates are the same in this and the other coords.
// It is assumed that this coord and other coord both have the same (provided)
// layout.
func (c Coord) Equal(layout Layout, other Coord) bool {

	numOrds := len(c)

	if layout.Stride() < numOrds {
		numOrds = layout.Stride()
	}

	if (len(c) < layout.Stride() || len(other) < layout.Stride()) && len(c) != len(other) {
		return false
	}

	for i := 0; i < numOrds; i++ {
		if math.IsNaN(c[i]) || math.IsNaN(other[i]) {
			if !math.IsNaN(c[i]) || !math.IsNaN(other[i]) {
				return false
			}
		} else if c[i] != other[i] {
			return false
		}
	}

	return true
}

// T is a generic interface implemented by all geometry types.
type T interface {
	// Layout defines how the coordinates are organized in the array of floats (returned by FlatCoords())
	Layout() Layout
	// Stride is the same as Layout().Stride()
	Stride() int
	// Bounds is the Bounds object that contains the geometry
	Bounds() *Bounds
	// FlatCoords contains the coordinates of the geometry packing into a single array of floats.
	// To interpret the coords, the Layout() must be used
	FlatCoords() []float64
	// Ends returns the ends of the sub geometries contained in this geometry.
	// If the geometry type does not support sub-geometries (like lines) then this
	// returns nil.
	// (MultiLine and Polygon support sub-geometries)
	Ends() []int
	// Endss returns nil unless this type is a MultiPolygon.  In the case of MultiPolygon
	// Endss returns an array of arrays where the inner array are essentially the Ends of the contained polygons
	Endss() [][]int
	// The Reference System code identifying a the reference system the geometry is encoded in.
	// The meaning and coding of the SRID is application dependent.
	SRID() int
	// Dimensionality of this type of geometry.  Defined in the OGC Simple Feature Specification
	// section 2.1.13.1.
	Dimensionality() ogc.Dimensionality
	// OGCBoundary returns the boundary geometry as defined in the OGC Simple Feature Specification
	// section 2.1.13.1.
	// If a geometry is empty (no coordinates) or otherwise does not have a boundary the result is
	// a MultiPoint with no points.
	// Points and Closed lines have no boundaries (returns MultiPoint containing no points)
	// Boundary of Lines are the end points of the lines
	// Boundary of a MultiCurve consists of those Points that are in the boundaries of an odd number of
	// its element Curves.
	// Boundary of a Polygon consists of its set of Rings.
	// Boundary of a MultiPolygon consists of the set of Rings of its Polygons.
	// Boundary of an arbitrary Collection of geometries whose interiors are disjoint consists of
	// geometries drawn from the boundaries of the element geometries by application of the ‘mod 2’ union rule
	OGCBoundary() T
	// OGCBoundaryDimensionality efficiently calculates of OGCBoundary().Dimensionality()
	// (Skips creating the geometry)
	OGCBoundaryDimensionality() ogc.Dimensionality
}

// MIndex returns the index of the M dimension, or -1 if the l does not have an
// M dimension.
func (l Layout) MIndex() int {
	switch l {
	case NoLayout, XY, XYZ:
		return -1
	case XYM:
		return 2
	case XYZM:
		return 3
	default:
		return 3
	}
}

// Stride returns l's number of dimensions.
func (l Layout) Stride() int {
	switch l {
	case NoLayout:
		return 0
	case XY:
		return 2
	case XYZ:
		return 3
	case XYM:
		return 3
	case XYZM:
		return 4
	default:
		return int(l)
	}
}

// String returns a human-readable string representing l.
func (l Layout) String() string {
	switch l {
	case NoLayout:
		return "NoLayout"
	case XY:
		return "XY"
	case XYZ:
		return "XYZ"
	case XYM:
		return "XYM"
	case XYZM:
		return "XYZM"
	default:
		return fmt.Sprintf("Layout(%d)", int(l))
	}
}

// ZIndex returns the index of l's Z dimension, or -1 if l does not have a Z
// dimension.
func (l Layout) ZIndex() int {
	switch l {
	case NoLayout, XY, XYM:
		return -1
	default:
		return 2
	}
}

// Must panics if err is not nil, otherwise it returns g.
func Must(g T, err error) T {
	if err != nil {
		panic(err)
	}
	return g
}

var (
	errIncorrectEnd         = errors.New("geom: incorrect end")
	errLengthStrideMismatch = errors.New("geom: length/stride mismatch")
	errMisalignedEnd        = errors.New("geom: misaligned end")
	errNonEmptyEnds         = errors.New("geom: non-empty ends")
	errNonEmptyEndss        = errors.New("geom: non-empty endss")
	errNonEmptyFlatCoords   = errors.New("geom: non-empty flatCoords")
	errOutOfOrderEnd        = errors.New("geom: out-of-order end")
	errStrideLayoutMismatch = errors.New("geom: stride/layout mismatch")
)
