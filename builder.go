package geom

import (
	"bytes"
	"fmt"
)

// Builder is a builder for creating arbitrary geometries.
// It performs validations during each step of the construction process
// to ensure that the geometries created are valid geometries once created.
// For example LinearRings will be closed, Polygon holes will be contained
// within the shell, etc...
type Builder struct {
	layout  Layout
	current []float64
	parts   []part
	err     error
}

// datastructures to use in builder while builder is in building phase
type part struct {
	geomType geomType
	data     []float64
}

type geomType int

const (
	unspecified geomType = iota
	point
	line
	polygon
)

func (t geomType) String() string {
	switch t {
	case unspecified:
		return "unspecified"
	case point:
		return "point"
	case line:
		return "line"
	case polygon:
		return "polygon"
	default:
		return "invalid geomType: " + string(t)
	}
}

func Build(layout Layout) *Builder {
	return &Builder{layout: layout}
}
func (b *Builder) GoTo(coord ...float64) *Builder {
	if b.err == nil {
		if len(coord) != b.layout.Stride() {
			b.err = fmt.Errorf("GoTo(%v) does not have the correct number of ordinates.  Layout indicates that %v are required", stringSlice(coord), b.layout.Stride())
			return b
		}
		if b.current != nil {
			b.parts = append(b.parts, part{unspecified, b.current})
		}
		b.current = coord
	}
	return b
}
func (b *Builder) LineTo(coord ...float64) *Builder {
	return b
}
func (b *Builder) AddLineSegments(coords ...Coord) *Builder {
	for _, c := range coords {
		b.LineTo(c...)
	}
	return b
}
func (b *Builder) CloseRing() *Builder {
	return b
}
func (b *Builder) Point() (point *Point, err error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.current == nil {
		return nil, fmt.Errorf("The GoTo method must first be called in order to provide the coordinate data before a point can be created")
	}
	if len(b.parts) != 0 {
		return nil, fmt.Errorf("The builder contains other parts beyond the coordinates of a point: %v", b.parts)
	}
	return NewPointFlat(b.layout, b.current), nil
}
func (b *Builder) MultiPoint() (points *MultiPoint, err error) {
	return nil, nil
}
func (b *Builder) LineString() (line *LineString, err error) {
	return nil, nil
}
func (b *Builder) LinearRing() (ring *LinearRing, err error) {
	return nil, nil
}
func (b *Builder) MultiLineString() (lines *MultiLineString, err error) {
	return nil, nil
}
func (b *Builder) Polygon() (poly *Polygon, err error) {
	return nil, nil
}
func (b *Builder) MultiPolygon() (polys *MultiPolygon, err error) {
	return nil, nil
}

func stringSlice(coord []float64) string {
	buf := bytes.Buffer{}
	for _, c := range coord {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprint(c))
	}
	return buf.String()
}
