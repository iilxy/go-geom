package boundary

// NodeRule implementations determine whether node points which are in boundaries of Lineal geometry components
// are in the boundary of the parent geometry collection.
//
// The SFS specifies a single kind of boundary node rule, the Mod2BoundaryNodeRule rule.
//
// However, other kinds of Boundary Node Rules are appropriate in specific situations
// (for instance, linear network topology usually follows the {@link EndPointBoundaryNodeRule}.)
//
// Some operations (such as Relate, Boundary and IsSimple) allow the BoundaryNodeRule to be specified,
// and respect the supplied rule when computing the results of the operation.
//
// An example use case for a non-SFS-standard Boundary Node Rule is that of checking that a set of LineStrings have
// valid linear network topology, when turn-arounds are represented as closed rings.
// In this situation, the entry road to the turn-around is only valid when it touches the turn-around ring
// at the single (common) endpoint.  This is equivalent to requiring the set of LineStrings to be
// simple under the EndPointBoundaryNodeRule.
//
// The SFS-standard Mod2BoundaryNodeRule is not sufficient to perform this test, since it
// states that closed rings have no boundary points.
type NodeRule interface {
	// IsInBoundary tests whether a point that lies in the boundaryCount
	// geometry component boundaries is considered to form part of the boundary
	// of the parent geometry.
	//
	// Param boundaryCount - the number of component boundaries that this point occurs in
	IsInBoundary(boundaryCount int) bool
}

// Mod2BoundaryNodeRule specifies that points are in the boundary of a
// lineal geometry iff the point lies on the boundary of an odd number of components.
// Under this rule LinearRings and closed LineStrings have an empty boundary.
//
// This is the rule specified by the OGC SFS and is the default rule used.
type Mod2BoundaryNodeRule struct{}

func (r Mod2BoundaryNodeRule) IsInBoundary(boundaryCount int) bool {
	return boundaryCount%2 == 1
}

// EndPointBoundaryNodeRule specifies that any points which are endpoints
// of lineal components are in the boundary of the parent geometry.
//
// This corresponds to the "intuitive" topological definition of boundary.
// Under this rule LinearRings have a non-empty boundary
// (the common endpoint of the underlying LineString).
//
// This rule is useful when dealing with linear networks.
//
// For example, it can be used to check whether linear networks are correctly noded.
// The usual network topology constraint is that linear segments may touch only at endpoints.
// In the case of a segment touching a closed segment (ring) at one point,
// the Mod2 rule cannot distinguish between the permitted case of touching at the
// node point and the invalid case of touching at some other interior (non-node) point.
// The EndPoint rule does distinguish between these cases, so is more appropriate for use.
type EndPointBoundaryNodeRule struct{}

func (r EndPointBoundaryNodeRule) IsInBoundary(boundaryCount int) bool {
	return boundaryCount > 0
}

// MultiValentEndPointBoundaryNodeRule determines that only endpoints with valency
// greater than 1 are on the boundary.
//
// This corresponds to the boundary of a MultiLineString being all the "attached"
// endpoints, but not the "unattached" ones.
type MultiValentEndPointBoundaryNodeRule struct{}

func (r MultiValentEndPointBoundaryNodeRule) IsInBoundary(boundaryCount int) bool {
	return boundaryCount > 1
}

// MonoValentEndPointBoundaryNodeRule determines that only endpoints
// with valency of exactly 1 are on the boundary.
//
// This corresponds to the boundary of a MultiLineString being all
// the "unattached" endpoints.
type MonoValentEndPointBoundaryNodeRule struct{}

func (r MonoValentEndPointBoundaryNodeRule) IsInBoundary(boundaryCount int) bool {
	return boundaryCount == 1
}
