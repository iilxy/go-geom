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
	ends     []int
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

func (b *Builder) StartPolygon(coord ...float64) *Builder {
	return b.startGeom("StartPolygon", polygonType, coord, true)
}

func (b *Builder) StartHole(coord ...float64) *Builder {
	b.CloseRing()
	b.validateNextCoord("StartHole", coord)
	b.current.data = append(b.current.data, coord...)
	return b
}

func (b *Builder) LineTo(coord ...float64) *Builder {
	b.validateNextCoord("LineTo", coord)
	if b.err != nil {
		return b
	}
	if b.current.partType == nullType {
		b.err = &BuilderError{"LineTo() cannot be executed without first starting a (non-point) geometry.", debug.Stack()}
		return b
	}

	if b.current.partType == pointType {
		b.err = &BuilderError{"LineTo() cannot be executed on a point geometry.", debug.Stack()}
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
	if b.err != nil {
		return b
	}
	if b.current.partType == nullType {
		b.err = &BuilderError{"A geometry must be started before it can be closed", debug.Stack()}
		return b
	}
	if len(b.current.data)/b.layout.Stride() < 3 {
		b.err = &BuilderError{"CloseRing() can only be called if the current geometry has at least 3 coordinates", debug.Stack()}
		return b
	}

	lastEnd := 0
	if len(b.current.ends) > 0 {
		lastEnd = b.current.ends[len(b.current.ends)-1]
	}

	firstCoord := Coord(b.current.data[lastEnd : lastEnd+b.layout.Stride()])
	lastCoord := Coord(b.current.data[len(b.current.data)-b.layout.Stride():])

	if !firstCoord.Equal(b.layout, lastCoord) {
		b.current.data = append(b.current.data, firstCoord...)
		end := len(b.current.data)
		b.current.ends = append(b.current.ends, end)
	}

	return b
}
func (b *Builder) Point() (point *Point, err *BuilderError) {
	b.validateGeomType("Point", pointType)
	if b.err != nil {
		return nil, b.err
	}
	return NewPointFlat(b.layout, b.current.data), nil
}
func (b *Builder) MultiPoint(allGeoms bool) (points *MultiPoint, err *BuilderError) {
	if b.err != nil {
		return nil, b.err
	}

	i := 0
	if !allGeoms {
		i = len(b.parts)
		for ; i >= 0 && b.parts[i-1].partType == pointType; i-- {
		}
	}

	coords := make([]float64, 0, len(b.parts)-i+1)

	for ; i < len(b.parts); i++ {
		if b.parts[i].partType != pointType {
			b.err = &BuilderError{fmt.Sprintf("The geometry at index %d was not a point it was a '%v'", i, b.parts[i].partType), debug.Stack()}
			return nil, b.err
		}
		coords = append(coords, b.parts[i].data...)
	}

	b.validateGeomType("Point", pointType)

	if b.err != nil {
		return nil, b.err
	}

	coords = append(coords, b.current.data...)
	return NewMultiPointFlat(b.layout, coords), nil
}
func (b *Builder) LineString() (line *LineString, err *BuilderError) {
	b.validateGeomType("LineString", lineType)
	if b.err != nil {
		return nil, b.err
	}
	return NewLineStringFlat(b.layout, b.current.data), nil
}
func (b *Builder) LinearRing() (ring *LinearRing, err *BuilderError) {
	b.validateGeomType("LinearRing", lineType)
	if b.err != nil {
		return nil, b.err
	}
	b.CloseRing()
	return NewLinearRingFlat(b.layout, b.current.data), nil
}
func (b *Builder) MultiLineString(allGeoms bool) (lines *MultiLineString, err *BuilderError) {
	if b.err != nil {
		return nil, b.err
	}

	i := 0
	if !allGeoms {
		i = len(b.parts)
		for ; i >= 1 && b.parts[i-1].partType == lineType; i-- {
		}
	}
	numGeoms := len(b.parts) - i + 1
	coords := make([]float64, 0, numGeoms*b.layout.Stride())
	ends := make([]int, 0, numGeoms)

	for ; i < len(b.parts); i++ {
		if b.parts[i].partType != lineType {
			b.err = &BuilderError{fmt.Sprintf("The geometry at index %d was not a line it was a '%v'", i, b.parts[i].partType), debug.Stack()}
			return nil, b.err
		}
		coords = append(coords, b.parts[i].data...)
		ends = append(ends, len(coords))
	}

	b.validateGeomType("Line", lineType)

	if b.err != nil {
		return nil, b.err
	}

	coords = append(coords, b.current.data...)
	ends = append(ends, len(coords))
	return NewMultiLineStringFlat(b.layout, coords, ends), nil
}
func (b *Builder) Polygon() (poly *Polygon, err *BuilderError) {
	b.validateGeomType("Polygon", polygonType)
	if b.err != nil {
		return nil, b.err
	}
	b.CloseRing()
	return NewPolygonFlat(b.layout, b.current.data, b.current.ends), nil
}
func (b *Builder) MultiPolygon(allGeoms bool) (polys *MultiPolygon, err *BuilderError) {
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
		b.current = part{partType, lastCoord, []int{}}
	} else {
		b.current = part{partType, coord, []int{}}
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

func (b *Builder) validateGeomType(humanReadable string, requiredType partType) {
	if b.err != nil {
		return
	}
	if b.current.partType != requiredType {
		b.err = &BuilderError{fmt.Sprintf("A %s cannot be created because the current geometry under construction is a '%v'.", humanReadable, b.current.partType), debug.Stack()}
	}
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

func validateGeomType(requiredType partType) {}
