package graph_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/internal/graph"
	"reflect"
	"testing"
)

func TestPlanarGraph_AddEdges(t *testing.T) {
	graph := graph.NewPlanarGraph(graph.DefaultNodeFactory{})
	edges := []*graph.Edge{
		graph.NewEdge(geom.XY, []float64{0, 0, 10, 0}, nil),
		graph.NewEdge(geom.XY, []float64{10, 0, 10, 10}, nil),
		graph.NewEdge(geom.XY, []float64{10, 10, 0, 10}, nil),
		graph.NewEdge(geom.XY, []float64{0, 10, 0, 0}, nil),
	}
	graph.AddEdges(edges)

	i := 0
	graph.WalkEdges(func(edge *graph.Edge) bool {
		if !reflect.DeepEqual(edge, edges[i]) {
			t.Fatalf("edges not added as expected %v %v", i, edge)
		}
		i++
		return true
	})

	if i != len(edges) {
		t.Errorf("Not all the edges were found in the graph: expected %v but was %v", len(edges), i)
	}
}
