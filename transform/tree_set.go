package transform

// TreeSet sorts the coordinates according to the Compare strategy and removes duplicates as
// dictated by the Equals function of the Compare strategy
type TreeSet struct {
	treeMap *TreeMap
}

// NewTreeSet creates a new TreeSet instance
func NewTreeSet(compare Compare) *CoordTreeSet {
	treeMap := NewTreeMap(compare)
	return &TreeSet{
		treeMap: treeMap,
	}
}

// Insert adds a new object to the tree set
// Returns true if the coordinate was added, false if it was already in the tree
func (set *TreeSet) Insert(obj interface{}) bool {
	return set.treeMap.Insert(obj, nil)
}

// Walk passes each element in the map to the visitor.  The order of visiting is from the element with the smallest key
// to the element with the largest key
func (set *TreeSet) Walk(visitor func(obj interface{})) {
	set.treeMap.Walk(func(key, value interface{}) {
		visitor(key)
	})
}
