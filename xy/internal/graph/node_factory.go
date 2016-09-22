package graph

import "github.com/twpayne/go-geom"

type NodeFactory interface {
	create(coord geom.Coord) *Node
}

type DefaultNodeFactory struct{}

var _ NodeFactory = DefaultNodeFactory{}

func (onf DefaultNodeFactory) create(coord geom.Coord) *Node {
	return &Node{
		edges: newEdgeEndStarCommon(),
		coord: coord}
}

type OverlayNodeFactory struct{}

var _ NodeFactory = OverlayNodeFactory{}

func (onf OverlayNodeFactory) create(coord geom.Coord) *Node {
	return &Node{
		coord: coord,
		edges: NewDirectedEdgeStar([]edgeEnd{}),
	}
}
