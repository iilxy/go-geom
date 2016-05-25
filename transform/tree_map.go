package transform

// Compare compares two coordinates for equality and magnitude
type Compare interface {
	IsEquals(o1, o2 interface{}) bool
	IsLess(o1, o2 interface{}) bool
}

type tree struct {
	left       *tree
	key, value interface{}
	right      *tree
}

// TreeSet sorts the coordinates according to the Compare strategy and removes duplicates as
// dictated by the Equals function of the Compare strategy
type TreeMap struct {
	compare Compare
	tree    *tree
	size    int
}

// NewTreeSet creates a new TreeSet instance
func NewTreeMap(compare Compare) *TreeMap {
	return &TreeMap{
		compare: compare,
	}
}

// Size returns the number of elements in the tree map
func (tm *TreeMap) Size() int {
	return tm.size
}

// Insert adds a new coordinate to the tree set
// the coordinate must be the same size as the Stride of the layout provided
// when constructing the TreeSet
// Returns true if the coordinate was added, false if it was already in the tree
func (set *TreeMap) Insert(key, value interface{}) bool {
	tree, added := set.insertImpl(set.tree, key, value)
	if added {
		set.tree = tree
		set.size++
	}

	return added
}

// Walk passes each element in the map to the visitor.  The order of visiting is from the element with the smallest key
// to the element with the largest key
func (tm *TreeMap) Walk(visitor func(key, value interface{})) {
	tm.walk(tm.tree, func(key, value interface{}) bool {
		visitor(key, value)
		return true
	})
}

// WalkInterruptible passes each element in the map to the visitor until false is returned from visitor.
// The order of visiting is from the element with the smallest key to the element with the largest key
func (tm *TreeMap) WalkInterruptible(visitor func(key, value interface{}) bool) {
	tm.walk(tm.tree, visitor)
}

func (tm *TreeMap) walk(t *tree, visitor func(key, value interface{}) bool) {
	if t == nil {
		return
	}
	tm.walk(t.left, visitor)
	visitor(t.key, t.value)
	tm.walk(t.right, visitor)
}

func (tm *TreeMap) insertImpl(t *tree, key, value interface{}) (*tree, bool) {
	if t == nil {
		return &tree{left: nil, key: key, value: value, right: nil}, true
	}

	if tm.compare.IsEquals(key, t.key) {
		return t, false
	}

	var added bool
	if tm.compare.IsLess(key, t.key) {
		t.left, added = tm.insertImpl(t.left, key, value)
	} else {
		t.right, added = tm.insertImpl(t.right, key, value)
	}

	return t, added
}
