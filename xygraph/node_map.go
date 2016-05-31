package xygraph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/sorting"
	"github.com/twpayne/go-geom/transform"
	"github.com/twpayne/go-geom/xy/location"
)

type coordCompare struct{}

var _ transform.Compare = coordCompare{}

func (c coordCompare) IsEquals(o1, o2 interface{}) bool {
	return o1.(geom.Coord).Equal(o2.(geom.Coord))
}
func (c coordCompare) IsLess(o1, o2 interface{}) bool {
	return sorting.IsLess2D([]float64(o1.(geom.Coord)), []float64(o2.(geom.Coord)))
}

type NodeFactory interface {
	create(coord geom.Coord) *Node
}

// NodeMap is map of nodes, indexed by the coordinate of the node
type NodeMap struct {
	nodeFactory NodeFactory
	nodeMap     *transform.TreeMap
}

func NewNodeMap(nodeFactory NodeFactory) NodeMap {
	return &NodeMap{
		nodeFactory: nodeFactory,
		nodeMap:     transform.NewTreeMap(coordCompare{}),
	}
}

func (nm *NodeMap) addCoordNode(coord geom.Coord) *Node {
	node, has := nm.nodeMap.Get(coord)
	if has {
		node = nm.nodeFactory.create(coord)
		nm.nodeMap.Insert(coord, node)
	}
	return node.(*Node)
}

func (nm *NodeMap) addNode(n Node) *Node {
	node, has := nm.nodeMap.Get(n.coord)
	if !has {
		nm.nodeMap.Insert(n.coord, n)
		return n
	}
	node.(*Node).mergeLabel(n)
	return node.(*Node)
}

// addEdgeEnd adds a node for the start point of this EdgeEnd (if one does not already exist in this map).
// Adds the EdgeEnd to the (possibly new) node.
func (nm *NodeMap) addEdgeEnd(e EdgeEnd) {
	p := e.Coordinate()
	n := nm.addNode(p)
	n.add(e)
}

func (nm *NodeMap) find(c geom.Coord) (*Node, bool) {
	n, has := nm.nodeMap.Get(c)
	return n.(*Node), has
}

func (nm *NodeMap) getBoundaryNodes(geomIndex int) []*Node {
	bdyNodes := make([]*Node, nm.nodeMap.Size(), 0)
	nm.nodeMap.Walk(func(key, value interface{}) {
		node := value.(*Node)
		if node.label[geomIndex][ON] == location.Boundary {
			bdyNodes = append(bdyNodes, node)
		}
	})

	return bdyNodes
}
