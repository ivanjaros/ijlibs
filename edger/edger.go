package edger

// Actions creates relationships for fixed-length bytes slice IDs.
// Each key consists of 1-byte long relationship type, 1-byte long relationship order
// and fixed-length byte slice ids of the left and right edges.
// All relationships are automatically bi-directional so parent can look for its children
// as well as children can look for their parents.
// Setting "root" can be done by defining relation for self, ie. the child and parent share the same id
// and then the rest of the records can point to this root as parent.

type Edge struct {
	Item  []byte
	Edges []Edge
}

// skipSelf allows to return only parents/children without id of self
func (e Edge) GetIds(skipSelf ...bool) [][]byte {
	var ids [][]byte

	if !(len(skipSelf) > 0 && skipSelf[0] == true) {
		ids = append(ids, e.Item)
	}

	for k := range e.Edges {
		ids = append(ids, e.Edges[k].GetIds()...)
	}

	return ids
}

type Actions interface {
	SaveEdges(relType byte, ltrPairs ...[]byte) error
	DeleteEdges(relType byte, ltrPairs ...[]byte) error
	// deletes all edges associated with provided item by deleting direct children and direct parents
	DeleteItem(relType byte, item []byte) error
	// loads items from the up/right direction
	LoadParents(relType byte, item []byte, maxLevel ...int) (Edge, error)
	// loads items from the down/left direction
	LoadChildren(relType byte, item []byte, maxLevel ...int) (Edge, error)
	// Loads all edges starting from child going up/right to the parent.
	LoadParentsUntil(relType byte, child, parent []byte) (Edge, error)
	// Loads all edges starting from parent down/left to the child.
	LoadChildrenUntil(relType byte, parent, child []byte) (Edge, error)
}

type Edger interface {
	Actions
}

type ACIDEdger interface {
	Transaction() Transaction
	Edger
}

type Transaction interface {
	Commit() error
	Rollback() error
	Actions
}
