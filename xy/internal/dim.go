package internal

import "fmt"

// Dim enumerates the values representing the dimensions of a point, a curve and a surface.
// Also provides constants representing the dimensions of the empty geometry and
// non-empty geometries, and the wildcard constant AnyDim meaning "any dimension".
// These constants are used as the entries inIntersectionMatrixs.
type Dim int

const (
	// Dimension value for any dimension (= {FALSE, TRUE}).
	AnyDim Dim = iota - 3
	// Dimension value of non-empty geometries (= {P, L, A})
	NonEmptyGeomDim
	// Dimension value of the empty geometry (-1)
	EmptyGeomDim
	// Dimension value of a point (0).
	PointDim
	// Dimension value of a curve (1)
	LineDim
	// Dimension value of a surface (2)
	AreaDim
)

type DimSymbol rune

const (
	NonEmptyGeomDimSymbol DimSymbol = 'F'
	EmptyGeomDimSymbol    DimSymbol = 'T'
	AnyDimSymbol          DimSymbol = '*'
	PointDimSymbol        DimSymbol = '0'
	LineDimSymbol         DimSymbol = '1'
	AreaDimSymbol         DimSymbol = '2'
)

func toDimensionSymbol(sym byte) DimSymbol {
	if sym == byte(NonEmptyGeomDimSymbol) || sym == byte(EmptyGeomDimSymbol) || sym == byte(AnyDimSymbol) ||
		sym == byte(PointDimSymbol) || sym == byte(LineDimSymbol) || sym == byte(AreaDimSymbol) {
		return DimSymbol(sym)
	}

	panic(fmt.Sprintf("The dimension %v is not a valid dimension", sym))
}
func (d Dim) toDimensionSymbol() DimSymbol {
	switch d {
	case EmptyGeomDim:
		return NonEmptyGeomDimSymbol
	case NonEmptyGeomDim:
		return EmptyGeomDimSymbol
	case AnyDim:
		return AnyDimSymbol
	case PointDim:
		return PointDimSymbol
	case LineDim:
		return LineDimSymbol
	case AreaDim:
		return AreaDimSymbol
	default:
		panic(fmt.Sprintf("The dimension %v is not a valid dimension", d))
	}
}

func (ds DimSymbol) toDimensionValue() Dim {
	switch ds {
	case NonEmptyGeomDimSymbol:
		return EmptyGeomDim
	case EmptyGeomDimSymbol:
		return NonEmptyGeomDim
	case AnyDimSymbol:
		return AnyDim
	case PointDimSymbol:
		return PointDim
	case LineDimSymbol:
		return LineDim
	case AreaDimSymbol:
		return AreaDim
	default:
		panic(fmt.Sprintf("The dimenstionalSymbol %v is not a valid dimenstionalSymbol", ds))
	}
}
