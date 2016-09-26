package graph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/location"
)

// T is the interface all graphs must satisfy
type T interface {
	Nodes() []Node
	Edges() []Edge
}

type PropertyName string

const (
	// Visited is a property indicating if the Graphable has been visited in the current algorithm
	// values are bool
	Visited PropertyName = "Visited"
	// Coord is a property which is a representative coordinate for the Graphable
	// values are geom.Coord
	Coord PropertyName = "Coord"
	// Isolated is a property indicating if the Graphable interacts or touches any other component.
	// This is the case if the Graphable only has a single RelatedObject
	// values are bool
	Isolated PropertyName = "Isolated"
)

// Graphable is the basic interface for all graph components (Nodes and Edge)
// Graphical objects can be associated/related to one or more geometries.  This
// can be used in certian algorithms like Buffer
type Graphable interface {
	// Related returns the associated/related geometries along with
	// how they are related to the current Graphable object
	Related() []RelatedObject
	// Properties is a storage area for algorithms to put data related to the algorithm
	// PropertyName constants are properties commonly used
	Properties() map[PropertyName]interface{}
}

// RelatedObject defines the relationship between to a geometry.
//
// If the related Graphable is a Node only Left and Right will be location.None
// and only On will have a valid location
//
// If the related Graphable is an Edge all three properties (Left, On, Right) will
// have valid locations
type RelatedObject struct {
	// Geom is the geometry that the Graphable is related to
	Geom geom.T
	// Left indicates where the geometry is with respect to the left of the Graphable.  Only applies if
	// the Graphable has a logical left (like Edge).  If not applicable this will be location.None.
	Left,
	// On indicates where the geometry is with respect to the Graphable.
	On,
	// Right indicates where the geometry is with respect to the right of the Graphable.  Only applies if
	// the Graphable has a logical right (like Edge)  If not applicable this will be location.None.
	Right location.Type
}

// GraphableImpl an implementation for the Graphable interface and can be used when implementing other Graphable
// objects (like Nodes and Edges)
type GraphableImpl struct {
	// Contain the related objects for retrieval in RelatedObject interface
	Related    []RelatedObject
	properties map[PropertyName]interface{}
}

var _ Graphable = GraphableImpl{}

func (g GraphableImpl) Related() []RelatedObject {
	return g.Related
}

func (g GraphableImpl) Properties() map[PropertyName]interface{} {
	return g.properties
}

// Node represents a node in a graph
type Node struct {
	GraphableImpl
}

var _ Graphable = Node{}
