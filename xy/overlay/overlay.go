package overlay

import "github.com/twpayne/go-geom/xy/location"

type Type interface {
	IsResultOfOp(loc1, loc2 location.Type) bool
}

func mapBoundary(loc1 location.Type) location.Type {
	if loc1 == location.Boundary {
		return location.Interior
	}
	return loc1
}

type Intersection struct{}
func (op Intersection) IsResultOfOp(loc0, loc1 location.Type) bool {
	return mapBoundary(loc0) == location.Interior && mapBoundary(loc1) == location.Interior;
}

type Union struct{}
func (op Union) IsResultOfOp(loc0, loc1 location.Type) bool {
	return mapBoundary(loc0) == location.Interior || mapBoundary(loc1) == location.Interior;
}

type Difference struct{}
func (op Difference) IsResultOfOp(loc0, loc1 location.Type) bool {
	return mapBoundary(loc0) == location.Interior && mapBoundary(loc1) != location.Interior;
}

type SymDifference struct{}
func (op SymDifference) IsResultOfOp(loc0, loc1 location.Type) bool {
	return (loc0 == location.Interior &&  loc1 != location.Interior) || (loc0 != location.Interior &&  loc1 == location.Interior);
}