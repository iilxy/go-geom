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

// Insert puts a key value pair into to the tree map.
// Returns true if a new entry was added, false if the key
// was already in the tree (the value may still have been updated)
func (set *TreeMap) Insert(key, value interface{}) bool {
	tree, added := set.insertImpl(set.tree, key, value)
	if added {
		set.tree = tree
		set.size++
	}

	return added
}

// Get returns the value associated with the key and true or nil and false
// the key does not have to be the same instance as the key in the map only
// be the equivalent key as evaluated by the Compare object configured in this tree
func (set *TreeMap) Get(key interface{}) (value interface{}, has bool) {
	return set.getImpl(set.tree, key)
}

// FindKey searches the key-set of this map for a matching key and returns the key instance from the map (if match found)
func (set *TreeMap) FindKey(key interface{}) (actual interface{}, has bool) {
	return set.findImpl(set.tree, key)
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

func (tm *TreeMap) walk(t *tree, visitor func(key, value interface{}) bool) bool {
	if t == nil {
		return true
	}
	if !tm.walk(t.left, visitor) {
		return false
	}
	if !visitor(t.key, t.value) {
		return false
	}
	return tm.walk(t.right, visitor)
}

func (tm *TreeMap) insertImpl(t *tree, key, value interface{}) (*tree, bool) {
	if t == nil {
		return &tree{left: nil, key: key, value: value, right: nil}, true
	}

	if tm.compare.IsEquals(key, t.key) {
		t.value = value
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

func (tm *TreeMap) getImpl(t *tree, key interface{}) (interface{}, bool) {
	switch {
	case t == nil:
		return nil, false
	case tm.compare.IsEquals(key, t.key):
		return t.value, true
	case tm.compare.IsLess(key, t.key):
		return tm.getImpl(t.left, key)
	default:
		return tm.getImpl(t.right, key)
	}
}

func (tm *TreeMap) findImpl(t *tree, key interface{}) (interface{}, bool) {
	switch {
	case t == nil:
		return nil, false
	case tm.compare.IsEquals(key, t.key):
		return t.key, true
	case tm.compare.IsLess(key, t.key):
		return tm.findImpl(t.left, key)
	default:
		return tm.findImpl(t.right, key)
	}
}
