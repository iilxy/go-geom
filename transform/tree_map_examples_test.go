package transform_test

import (
	"fmt"
	"github.com/twpayne/go-geom/transform"
)

type exampleCompare struct{}

func (c exampleCompare) IsEquals(o1, o2 interface{}) bool {
	i1, i2 := o1.(int), o2.(int)
	return i1 == i2
}
func (c exampleCompare) IsLess(o1, o2 interface{}) bool {
	i1, i2 := o1.(int), o2.(int)
	return i1 < i2
}

func ExampleNewTreeMap() {
	treeMap := transform.NewTreeMap(exampleCompare{})
	treeMap.Insert(1, "One")
	treeMap.Insert(3, "Three")
	treeMap.Insert(2, "Two")
	treeMap.Insert(1, "One")

	fmt.Printf("Size: %v Elements: ", treeMap.Size())

	treeMap.Walk(func(k, v interface{}) {
		fmt.Printf("%v-%v, ", k, v)
	})

	// Output: Size: 3 Elements: 1-One, 2-Two, 3-Three,
}
