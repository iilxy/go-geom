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

// Find searches the set for an equivalent object and returns the object or nil if the object is not within the set
func (set *TreeSet) Find(obj interface{}) (actual interface{}, has bool) {
	return set.Find(obj)
}

// Walk passes each element in the set to the visitor.  The order of visiting is from the element with the smallest key
// to the element with the largest key
func (set *TreeSet) Walk(visitor func(obj interface{})) {
	set.treeMap.Walk(func(key, value interface{}) {
		visitor(key)
	})
}

// WalkInterruptible passes each element in the set to the visitor until false is returned from visitor.
// The order of visiting is from the element with the smallest value to the element with the largest key
func (set *TreeSet) WalkInterruptible(visitor func(obj interface{}) bool) {
	set.treeMap.WalkInterruptible(func(key, value interface{}) bool {
		return visitor(key)
	})
}
