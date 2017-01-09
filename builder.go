package geom

import (
	"bytes"
	"fmt"
	"runtime/debug"
)

// Builder is a builder for creating arbitrary geometries.
// It performs validations during each step of the construction process
// to ensure that the geometries created are valid geometries once created.
// For example LinearRings will be closed, Polygon holes will be contained
// within the shell, etc...
type Builder struct {
	layout  Layout
	current part
	parts   []part
	err     *BuilderError
}

type BuilderError struct {
	Msg   string
	Stack []byte
}

func (err *BuilderError) String() string {
	return fmt.Sprintf("%s: \nStackTrace: %s", err.Msg, string(err.Stack))
}

// datastructures to use in builder while builder is in building phase
type part struct {
	partType partType
	data     []float64
}

type partType int

const (
	nullType partType = iota
	pointType
	lineType
	polygonType
)

func (t partType) String() string {
	switch t {
	case nullType:
		return "null type"
	case pointType:
		return "point"
	case lineType:
		return "line"
	case polygonType:
		return "polygon"
	default:
		return "invalid partType: " + string(t)
	}
}

func Build(layout Layout) *Builder {
	return &Builder{layout: layout}
}
func (b *Builder) StartPoint(coord ...float64) *Builder {
	return b.startGeom("StartPoint", pointType, coord, false)
}

func (b *Builder) StartLine(coord ...float64) *Builder {
	return b.startGeom("StartLine", lineType, coord, true)
}

func (b *Builder) LineTo(coord ...float64) *Builder {
	b.validateNextCoord("LineTo", coord)
	if b.err != nil {
		return b
	}
	if b.current.partType == nullType {
		b.err = &BuilderError{"LineTo() cannot be executed without first starting a (non-point) geometry.", debug.Stack()}
	}

	if b.current.partType == pointType {
		b.err = &BuilderError{"LineTo() cannot be executed on a point geometry.", debug.Stack()}
	}

	if b.err != nil {
		return b
	}

	b.current.data = append(b.current.data, coord...)

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
func (b *Builder) Point() (point *Point, err *BuilderError) {
	if b.err != nil {
		return nil, b.err
	}
	if b.current.partType != pointType {
		return nil, &BuilderError{fmt.Sprintf("A point cannot be created because the current geometry under construction is a '%v'.", b.current.partType), debug.Stack()}
	}
	return NewPointFlat(b.layout, b.current.data), nil
}
func (b *Builder) MultiPoint() (points *MultiPoint, err *BuilderError) {
	return nil, nil
}
func (b *Builder) LineString() (line *LineString, err *BuilderError) {
	if b.err != nil {
		return nil, b.err
	}

	if b.current.partType != lineType {
		return nil, &BuilderError{fmt.Sprintf("A LineString cannot be created because the current geometry under construction is a '%v'.", b.current.partType), debug.Stack()}
	}

	return NewLineStringFlat(b.layout, b.current.data), nil
}
func (b *Builder) LinearRing() (ring *LinearRing, err *BuilderError) {
	return nil, nil
}
func (b *Builder) MultiLineString() (lines *MultiLineString, err *BuilderError) {
	return nil, nil
}
func (b *Builder) Polygon() (poly *Polygon, err *BuilderError) {
	return nil, nil
}
func (b *Builder) MultiPolygon() (polys *MultiPolygon, err *BuilderError) {
	return nil, nil
}

func (b *Builder) validateNextCoord(methodName string, coord []float64) {
	if b.err == nil {
		if len(coord) != b.layout.Stride() {
			b.err = &BuilderError{fmt.Sprintf("%s(%v) does not have the correct number of ordinates.  Layout indicates that %v are required", methodName, stringSlice(coord), b.layout.Stride()), debug.Stack()}
		}
	}
}
func (b *Builder) startGeom(methodName string, partType partType, coord []float64, continueFromLastGeom bool) *Builder {
	if !continueFromLastGeom || len(coord) != 0 {
		b.validateNextCoord(methodName, coord)
	} else if len(b.parts) == 0 || len(b.parts[len(b.parts)-1].data) == 0 {
		b.err = &BuilderError{fmt.Sprintf("%v() can only be used when a geometry has already been created", methodName), debug.Stack()}
	}
	b.endGeom()

	if b.err != nil {
		return b
	}
	if continueFromLastGeom && len(coord) == 0 {
		lastPartCoords := b.parts[len(b.parts)-1].data
		lastCoord := lastPartCoords[len(lastPartCoords)-b.layout.Stride():]
		b.current = part{partType, lastCoord}
	} else {
		b.current = part{partType, coord}
	}
	return b
}
func (b *Builder) endGeom() {
	if b.err != nil {
		return
	}

	switch b.current.partType {
	case nullType:
		b.current = part{partType: nullType}
	case polygonType:
	//todo
	default:
		b.parts = append(b.parts, b.current)
	}
	b.current = part{partType: nullType}
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

func formatParts(parts []part) string {
	buf := bytes.Buffer{}
	for _, p := range parts {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("{type: \"%v\", coordinates: (%s)}", p.partType, stringSlice(p.data)))
	}
	return buf.String()
}
