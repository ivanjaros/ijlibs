package medger

import (
	"bytes"
	"errors"
	"github.com/ivanjaros/ijlibs/edger"
	edger_keys "github.com/ivanjaros/ijlibs/edger/keys"
	"github.com/plar/go-adaptive-radix-tree"
	"sync"
)

func New(keySize int, prefix ...byte) (*mEdger, error) {
	if keySize < 1 {
		return nil, errors.New("invalid key length")
	}
	return &mEdger{tree: art.New(), px: prefix, kl: keySize}, nil
}

type mEdger struct {
	tree art.Tree
	px   []byte
	kl   int
	mx   sync.RWMutex
}

func (m *mEdger) SaveEdges(relType byte, ltrPairs ...[]byte) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	keys, err := edger_keys.BuildEdgesForPairs(relType, m.kl, m.px, ltrPairs...)
	if err != nil {
		return err
	}

	for _, key := range keys {
		m.tree.Insert(key, nil)
	}

	return nil
}

func (m *mEdger) DeleteEdges(relType byte, ltrPairs ...[]byte) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	keys, err := edger_keys.BuildEdgesForPairs(relType, m.kl, m.px, ltrPairs...)
	if err != nil {
		return err
	}

	for _, key := range keys {
		m.tree.Delete(key)
	}

	return nil
}

func (m *mEdger) DeleteItem(relType byte, item []byte) error {
	parents, err := m.LoadParents(relType, item, 1)
	if err != nil {
		return err
	}

	children, err := m.LoadChildren(relType, item, 1)
	if err != nil {
		return err
	}

	var pairs [][]byte
	for k := range parents.Edges {
		pairs = append(pairs, item, parents.Edges[k].Item)
	}
	for k := range children.Edges {
		pairs = append(pairs, children.Edges[k].Item, item)
	}

	return m.DeleteEdges(relType, pairs...)
}

func (m *mEdger) LoadParentsUntil(relType byte, child, parent []byte) (edger.Edge, error) {
	return m.loadRange(relType, edger_keys.Ltr, child, parent)
}

func (m *mEdger) LoadChildrenUntil(relType byte, parent, child []byte) (edger.Edge, error) {
	return m.loadRange(relType, edger_keys.Rtl, parent, child)
}

func (m *mEdger) loadRange(relType byte, way byte, left, right []byte) (edger.Edge, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	result := edger.Edge{Item: left}
	rangeReader(m.tree, &result, relType, way, right, m.px...)
	return result, nil
}

func (m *mEdger) LoadParents(relType byte, item []byte, maxLevel ...int) (edger.Edge, error) {
	return m.loadRelations(relType, edger_keys.Ltr, item, maxLevel...)
}

func (m *mEdger) LoadChildren(relType byte, item []byte, maxLevel ...int) (edger.Edge, error) {
	return m.loadRelations(relType, edger_keys.Rtl, item, maxLevel...)
}

func (m *mEdger) loadRelations(relType byte, way byte, child []byte, maxLevel ...int) (edger.Edge, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	result := edger.Edge{Item: child}

	max := -1
	if len(maxLevel) > 0 && maxLevel[0] > 0 {
		max = maxLevel[0]
	}

	loopReader(m.tree, &result, relType, way, max, m.px...)

	return result, nil
}

func rangeReader(tree art.Tree, rel *edger.Edge, relType byte, relOrd byte, until []byte, p ...byte) {
	prefix := append(p, edger_keys.NewPrefix(relType, relOrd, rel.Item)...)
	found := getPrefixed(tree, prefix)
	var roots [][]byte
	var done bool
	for k := range found {
		left, right := edger_keys.GetEdges(found[k][len(p):])
		rel.Edges = append(rel.Edges, edger.Edge{Item: right})
		if bytes.Equal(left, right) {
			roots = append(roots, right)
		}
		if bytes.Equal(right, until) || bytes.Equal(left, until) {
			done = true
		}
	}

	if done {
		return
	}

	for k := range rel.Edges {
		for _, root := range roots {
			if bytes.Equal(rel.Edges[k].Item, root) {
				continue
			}
		}
		rangeReader(tree, &rel.Edges[k], relType, relOrd, until, p...)
	}
}

func loopReader(tree art.Tree, rel *edger.Edge, relType byte, relOrd byte, maxLoops int, p ...byte) {
	prefix := append(p, edger_keys.NewPrefix(relType, relOrd, rel.Item)...)
	found := getPrefixed(tree, prefix)
	var roots [][]byte
	for k := range found {
		left, right := edger_keys.GetEdges(found[k][len(p):])
		rel.Edges = append(rel.Edges, edger.Edge{Item: right})
		if bytes.Equal(left, right) {
			roots = append(roots, right)
		}
	}

	if maxLoops == 0 {
		return
	}

	for k := range rel.Edges {
		for _, root := range roots {
			if bytes.Equal(rel.Edges[k].Item, root) {
				continue
			}
		}
		loopReader(tree, &rel.Edges[k], relType, relOrd, maxLoops-1, p...)
	}
}

func getPrefixed(tree art.Tree, prefix []byte) [][]byte {
	var matches [][]byte

	tree.ForEachPrefix(prefix, func(node art.Node) bool {
		if key := node.Key(); key != nil {
			cp := make([]byte, len(key))
			copy(cp, key)
			matches = append(matches, key)
		}
		return true
	})

	return matches
}
