package xygraph

import (
	"fmt"
	"github.com/twpayne/go-geom/xy/location"
)

// IntersectionMatrix models a <b>Dimensionally Extended Nine-Intersection Model (DE-9IM)</b> matrix.
// DE-9IM matrices (such as "212FF1FF2")
// specify the topological relationship between two {@link Geometry}s.
// This class can also represent matrix patterns (such as "T*T******")
// which are used for matching instances of DE-9IM matrices.
//
//  Methods are provided to:
//  <UL>
//    <LI> set and query the elements of the matrix in a convenient fashion
//    <LI> convert to and from the standard string representation (specified in
//    SFS Section 2.1.13.2).
//    <LI> test to see if a matrix matches a given pattern string.
//  </UL>
//  <P>
//
//  For a description of the DE-9IM and the spatial predicates derived from it,
//  see the <i><A
//  HREF="http://www.opengis.org/techno/specs.htm">OGC 99-049 OpenGIS Simple Features
//  Specification for SQL</A></i>, as well as
//  <i>OGC 06-103r4 OpenGIS
//  Implementation Standard for Geographic information -
//  Simple feature access - Part 1: Common architecture</i>
//  (which provides some further details on certain predicate specifications).
// <p>
// The entries of the matrix are defined by the constants in the {@link Dimension} class.
// The indices of the matrix represent the topological locations
// that occur in a geometry (Interior, Boundary, Exterior).
// These are provided as constants in the {@link Location} class.
type IntersectionMatrix [3][3]dimension

func NewNullIntersectionMatrix() *IntersectionMatrix {
	return &IntersectionMatrix{[3]dimension{dimFALSE, dimFALSE, dimFALSE}, [3]dimension{dimFALSE, dimFALSE, dimFALSE}}
}

func NewIntersectionMatrixFromTemplate(template IntersectionMatrix) *IntersectionMatrix {
	im := &IntersectionMatrix{}
	im[0][0] = template[0][0]
	im[0][1] = template[0][1]
	im[0][2] = template[0][2]

	im[1][0] = template[1][0]
	im[1][1] = template[1][1]
	im[1][2] = template[1][2]

	im[2][0] = template[2][0]
	im[2][1] = template[2][1]
	im[2][2] = template[2][2]
	return im
}

// Add adds one matrix to another.
// Addition is defined by taking the maximum dimension value of each position
// in the summand matrices.
func (im *IntersectionMatrix) Add(source IntersectionMatrix) {
	for i, a := range source {
		for j, v := range a {
			im.SetAtLeast(i, j, v)
		}
	}
}

// IsTrue tests if the dimension value matches dimension_TRUE (i.e.  has value 0, 1, 2 or TRUE).
func (im *IntersectionMatrix) IsTrue(actualDimensionValue dimension) bool {
	if actualDimensionValue >= 0 || actualDimensionValue == dimTRUE {
		return true
	}
	return false
}

// Matches tests if the dimension value satisfies the dimension symbol
func (im *IntersectionMatrix) Matches(actualDimensionValue dimension, requiredDimensionSymbol dimensionalSymbol) bool {
	switch {
	case requiredDimensionSymbol == SYM_DONTCARE:
		return true
	case requiredDimensionSymbol == SYM_TRUE && (actualDimensionValue >= 0 || actualDimensionValue == dimTRUE):
		return true
	case requiredDimensionSymbol == SYM_FALSE && actualDimensionValue == dimFALSE:
		return true
	case requiredDimensionSymbol == SYM_P && actualDimensionValue == dimP:
		return true
	case requiredDimensionSymbol == SYM_L && actualDimensionValue == dimL:
		return true
	case requiredDimensionSymbol == SYM_A && actualDimensionValue == dimA:
		return true
	default:
		return false
	}
}

// Set changes the elements of this IntersectionMatrix to the dimension symbols in dimensionSymbols.
// Param dimensionSymbols - nine dimension symbols to which to set this IntersectionMatrix s elements.
// Possible values are T, F, * , 0, 1, 2
func (im *IntersectionMatrix) Set(dimensionSymbols string) {
	for i, sym := range dimensionSymbols {
		row := i / 3
		col := i % 3
		im[row][col] = dimensionalSymbol(sym).toDimensionValue()
	}
}

// SetAtLeast changes the specified element to minimumDimensionValue if the element is less.
func (im *IntersectionMatrix) SetAtLeast(row, column int, minimumDimensionValue dimension) {
	if im[row][column] < minimumDimensionValue {
		im[row][column] = minimumDimensionValue
	}
}

// SetAtLeastIfValid changes the specified element to minimumDimensionValue if the element is less.
// Does nothing if row < 0 or column < 0.
func (im *IntersectionMatrix) SetAtLeastIfValid(row, column int, minimumDimensionValue dimension) {
	if row >= 0 && column >= 0 {
		im.SetAtLeast(row, column, minimumDimensionValue)
	}
}

// SetAtLeastFromSymbols changes the element to the corresponding minimum dimension symbol if the element
// is less for each element in this IntersectionMatrix
func (im *IntersectionMatrix) SetAtLeastFromSymbols(minimumDimensionSymbols string) {
	for i, sym := range minimumDimensionSymbols {
		row := i / 3
		col := i % 3
		im.SetAtLeast(row, col, dimensionalSymbol(sym).toDimensionValue())
	}
}

// SetAll changes the elements of this IntersectionMatrix to dimensionValue
func (im *IntersectionMatrix) SetAll(dimensionValue dimension) {
	for ai := 0; ai < 3; ai++ {
		for bi := 0; bi < 3; bi++ {
			im[ai][bi] = dimensionValue
		}
	}
}

// Disjoint Returns true if this IntersectionMatrix is *  FF*FF****. (no itersections)
func (im *IntersectionMatrix) Disjoint() bool {
	return im[location.Interior][location.Interior] == dimFALSE &&
		im[location.Interior][location.Boundary] == dimFALSE &&
		im[location.Boundary][location.Interior] == dimFALSE &&
		im[location.Boundary][location.Boundary] == dimFALSE
}

// Touches returns true if this IntersectionMatrix is FT*******, F**T***** or F***T****
func (im *IntersectionMatrix) Touches(dimensionOfGeometryA, dimensionOfGeometryB dimension) bool {
	if dimensionOfGeometryA > dimensionOfGeometryB {
		//no need to get transpose because pattern matrix is symmetrical
		return im.Touches(dimensionOfGeometryB, dimensionOfGeometryA)
	}
	if (dimensionOfGeometryA == dimA && dimensionOfGeometryB == dimA) ||
		(dimensionOfGeometryA == dimL && dimensionOfGeometryB == dimL) ||
		(dimensionOfGeometryA == dimL && dimensionOfGeometryB == dimA) ||
		(dimensionOfGeometryA == dimP && dimensionOfGeometryB == dimA) ||
		(dimensionOfGeometryA == dimP && dimensionOfGeometryB == dimL) {
		return im[location.Interior][location.Interior] == dimFALSE && (im.IsTrue(im[location.Interior][location.Boundary]) ||
			im.IsTrue(im[location.Boundary][location.Interior]) || im.IsTrue(im[location.Boundary][location.Boundary]))
	}

	return false
}

// Crosses tests whether this geometry crosses the specified geometry.
//
// The crosses< predicate has the following equivalent definitions:
//
// * The geometries have some but not all interior points in common.
// * The DE-9IM Intersection Matrix for the two geometries is
//   * T*T****** (for P/L, P/A, and L/A situations)
//   * T*****T** (for L/P, L/A, and A/L situations)
//   * 0******** (for L/L situations)

// For any other combination of dimensions this predicate returns false.
//
// The SFS defined this predicate only for P/L, P/A, L/L, and L/A situations.
// JTS extends the definition to apply to L/P, A/P and A/L situations as well.
// This makes the relation symmetric.
func (im *IntersectionMatrix) Crosses(dimensionOfGeometryA, dimensionOfGeometryB dimension) bool {
	if (dimensionOfGeometryA == dimP && dimensionOfGeometryB == dimL) ||
		(dimensionOfGeometryA == dimP && dimensionOfGeometryB == dimA) ||
		(dimensionOfGeometryA == dimL && dimensionOfGeometryB == dimA) {
		return im.IsTrue(im[location.Interior][location.Interior]) &&
			im.IsTrue(im[location.Interior][location.Exterior])
	}
	if (dimensionOfGeometryA == dimL && dimensionOfGeometryB == dimP) ||
		(dimensionOfGeometryA == dimA && dimensionOfGeometryB == dimP) ||
		(dimensionOfGeometryA == dimA && dimensionOfGeometryB == dimL) {
		return im.IsTrue(im[location.Interior][location.Interior]) &&
			im.IsTrue(im[location.Exterior][location.Interior])
	}
	if dimensionOfGeometryA == dimL && dimensionOfGeometryB == dimL {
		return im[location.Interior][location.Interior] == 0
	}
	return false
}

// Within tests whether this IntersectionMatrix is T*F**F***
func (im *IntersectionMatrix) Within(dimensionOfGeometryA, dimensionOfGeometryB int) bool {
	return im.IsTrue(im[location.Interior][location.Interior]) &&
		im[location.Interior][location.Exterior] == dimFALSE &&
		im[location.Boundary][location.Exterior] == dimFALSE
}

// Contains tests whether this IntersectionMatrix is  T*****FF*
func (im *IntersectionMatrix) Contains(dimensionOfGeometryA, dimensionOfGeometryB int) bool {
	return im.IsTrue(im[location.Interior][location.Interior]) &&
		im[location.Exterior][location.Interior] == dimFALSE &&
		im[location.Exterior][location.Boundary] == dimFALSE
}

// Covers tests if this IntersectionMatrix is:
// * T*****FF*
// * or *T****FF*
// * or ***T**FF*
// * or ****T*FF*
func (im *IntersectionMatrix) Covers(dimensionOfGeometryA, dimensionOfGeometryB int) bool {
	hasPointInCommon := im.IsTrue(im[location.Interior][location.Interior]) ||
		im.IsTrue(im[location.Interior][location.Boundary]) ||
		im.IsTrue(im[location.Boundary][location.Interior]) ||
		im.IsTrue(im[location.Boundary][location.Boundary])

	return hasPointInCommon &&
		im[location.Exterior][location.Interior] == dimFALSE &&
		im[location.Exterior][location.Boundary] == dimFALSE
}

// CoveredBy tests if this IntersectionMatrix is
//  * T*F**F***
//  * or *TF**F***
//  * or **FT*F***
//  * or **F*TF***
func (im *IntersectionMatrix) CoveredBy(dimensionOfGeometryA, dimensionOfGeometryB int) bool {
	hasPointInCommon := im.IsTrue(im[location.Interior][location.Interior]) ||
		im.IsTrue(im[location.Interior][location.Boundary]) ||
		im.IsTrue(im[location.Boundary][location.Interior]) ||
		im.IsTrue(im[location.Boundary][location.Boundary])

	return hasPointInCommon &&
		im[location.Interior][location.Exterior] == dimFALSE &&
		im[location.Boundary][location.Exterior] == dimFALSE
}

// Equal tests whether the argument dimensions are equal and this IntersectionMatrix matches the pattern T*F**FFF*
//
// Note: This pattern differs from the one stated in Simple feature access - Part 1: Common architecture
// That document states the pattern as TFFFTFFFT.  This would specify that
// two identical POINTs are not equal, which is not desirable behaviour.
// The pattern used here has been corrected to compute equality in this situation.
func (im *IntersectionMatrix) Equal(dimensionOfGeometryA, dimensionOfGeometryB int) bool {
	if dimensionOfGeometryA != dimensionOfGeometryB {
		return false
	}
	return im.IsTrue(im[location.Interior][location.Interior]) &&
		im[location.Interior][location.Exterior] == dimFALSE &&
		im[location.Boundary][location.Exterior] == dimFALSE &&
		im[location.Exterior][location.Interior] == dimFALSE &&
		im[location.Exterior][location.Boundary] == dimFALSE
}

// Overlaps tests if this IntersectionMatrix is
// * T*T***T** (for two points or two surfaces)
// * 1*T***T** (for two curves)
func (im *IntersectionMatrix) Overlaps(dimensionOfGeometryA, dimensionOfGeometryB dimension) bool {
	if (dimensionOfGeometryA == dimP && dimensionOfGeometryB == dimP) ||
		(dimensionOfGeometryA == dimA && dimensionOfGeometryB == dimA) {
		return im.IsTrue(im[location.Interior][location.Interior]) &&
			im.IsTrue(im[location.Interior][location.Exterior]) &&
			im.IsTrue(im[location.Exterior][location.Interior])
	}
	if dimensionOfGeometryA == dimL && dimensionOfGeometryB == dimL {
		return im[location.Interior][location.Interior] == 1 &&
			im.IsTrue(im[location.Interior][location.Exterior]) &&
			im.IsTrue(im[location.Exterior][location.Interior])
	}
	return false
}

// MatchesSymbols tests whether the elements of this IntersectionMatrix satisfies the required dimension symbols.
func (im *IntersectionMatrix) MatchesSymbols(requiredDimensionSymbols string) bool {
	for ai := 0; ai < 3; ai++ {
		for bi := 0; bi < 3; bi++ {
			if !im.Matches(im[ai][bi], toDimensionSymbol(requiredDimensionSymbols[3*ai+bi])) {
				return false
			}
		}
	}
	return true
}

// Transpose transposes this IntersectionMatrix
func (im *IntersectionMatrix) Transpose() *IntersectionMatrix {
	im[1][0], im[0][1] = im[0][1], im[1][0]
	im[2][0], im[0][2] = im[0][2], im[2][0]
	im[2][1], im[1][2] = im[1][2], im[2][1]
	return im
}

// String Returns a nine-character String representation of this IntersectionMatrix
func (im IntersectionMatrix) String() string {
	buf := []byte("123456789")
	for ai := 0; ai < 3; ai++ {
		for bi := 0; bi < 3; bi++ {
			buf[3*ai+bi] = byte(im[ai][bi].toDimensionSymbol())
		}
	}
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7], buf[8])
}
