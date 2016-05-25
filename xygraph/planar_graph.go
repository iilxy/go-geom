package xygraph

// PlanarGraph contains nodes and edges corresponding to the nodes and line segments of
// a Geometry. Each node and edge in the graph is labeled with its topological location
// relative to the source geometry.
//
// Note that there is no requirement that points of self-intersection be a vertex.
// Thus to obtain a correct topology graph, Geometrys must be
// self-noded before constructing their graphs.
//
// Two fundamental operations are supported by topology graphs:
// * Computing the intersections between all the edges and nodes of a single graph
// * Computing the intersections between the edges and nodes of two different graphs
type PlanarGraph struct {
}
