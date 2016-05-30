package xygraph

import "fmt"

// dimension enumerates the values representing the dimensions of a point, a curve and a surface.
// Also provides constants representing the dimensions of the empty geometry and
// non-empty geometries, and the wildcard constant {@link #DONTCARE} meaning "any dimension".
// These constants are used as the entries in {@link IntersectionMatrix}s.
type dimension int

const (
	// Dimension value for any dimension (= {FALSE, TRUE}).
	dimDONTCARE dimension = iota - 3
	// Dimension value of non-empty geometries (= {P, L, A})
	dimTRUE
	// Dimension value of the empty geometry (-1)
	dimFALSE
	// Dimension value of a point (0).
	dimP
	// Dimension value of a curve (1)
	dimL
	// Dimension value of a surface (2)
	dimA
)

type dimensionalSymbol rune

const (
	SYM_FALSE    dimensionalSymbol = 'F'
	SYM_TRUE     dimensionalSymbol = 'T'
	SYM_DONTCARE dimensionalSymbol = '*'
	SYM_P        dimensionalSymbol = '0'
	SYM_L        dimensionalSymbol = '1'
	SYM_A        dimensionalSymbol = '2'
)

func toDimensionSymbol(sym byte) dimensionalSymbol {
	if sym == byte(SYM_FALSE) || sym == byte(SYM_TRUE) || sym == byte(SYM_DONTCARE) ||
		sym == byte(SYM_P) || sym == byte(SYM_L) || sym == byte(SYM_A) {
		return dimensionalSymbol(sym)
	}

	panic(fmt.Sprintf("The dimension %v is not a valid dimension", sym))
}
func (d dimension) toDimensionSymbol() dimensionalSymbol {
	switch d {
	case dimFALSE:
		return SYM_FALSE
	case dimTRUE:
		return SYM_TRUE
	case dimDONTCARE:
		return SYM_DONTCARE
	case dimP:
		return SYM_P
	case dimL:
		return SYM_L
	case dimA:
		return SYM_A
	default:
		panic(fmt.Sprintf("The dimension %v is not a valid dimension", d))
	}
}

func (ds dimensionalSymbol) toDimensionValue() dimension {
	switch ds {
	case SYM_FALSE:
		return dimFALSE
	case SYM_TRUE:
		return dimTRUE
	case SYM_DONTCARE:
		return dimDONTCARE
	case SYM_P:
		return dimP
	case SYM_L:
		return dimL
	case SYM_A:
		return dimA
	default:
		panic(fmt.Sprintf("The dimenstionalSymbol %v is not a valid dimenstionalSymbol", ds))
	}
}
