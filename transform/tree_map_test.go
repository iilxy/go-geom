package transform_test

import (
	"github.com/twpayne/go-geom/transform"
	"testing"
)

func TestTreeMap_Insert(t *testing.T) {
	treeMap := transform.NewTreeMap(exampleCompare{})
	if !treeMap.Insert(1, "_1_") {
		t.Fatalf("Insert did not report as being inserted")
	}
	if treeMap.Size() != 1 {
		t.Fatalf("Size was not 1 as expected, was: %v", treeMap.Size())
	}

	if !treeMap.Insert(3, "Three") {
		t.Fatalf("Insert did not report as being inserted")
	}
	if treeMap.Size() != 2 {
		t.Fatalf("Size was not 2 as expected, was: %v", treeMap.Size())
	}

	if !treeMap.Insert(2, "Two") {
		t.Fatalf("Insert did not report as being inserted")
	}
	if treeMap.Size() != 3 {
		t.Fatalf("Size was not 3 as expected, was: %v", treeMap.Size())
	}

	if treeMap.Insert(1, "One") {
		t.Fatalf("treeMap.Insert(1, \"One\") reported as being added but shouldn't have since key 1 already existed")
	}
	if treeMap.Size() != 3 {
		t.Fatalf("Size was not 3 as expected, was: %v", treeMap.Size())
	}
}

func TestTreeMap_Get(t *testing.T) {

	treeMap := transform.NewTreeMap(exampleCompare{})
	treeMap.Insert(1, "_1_")
	treeMap.Insert(3, "Three")
	treeMap.Insert(2, "Two")
	treeMap.Insert(1, "One")

	for _, tc := range []struct {
		idx      int
		has      bool
		expected interface{}
	}{
		{idx: 1, has: true, expected: "One"},
		{idx: 2, has: true, expected: "Two"},
		{idx: 3, has: true, expected: "Three"},
		{idx: 4, has: false, expected: nil},
	} {
		value, has := treeMap.Get(tc.idx)

		if has != tc.has {
			t.Errorf("Expected has to be %v but was %v", tc.has, has)
		}

		if value != tc.expected {
			t.Errorf("Expected to get %v but was %v", tc.expected, value)
		}
	}
}
