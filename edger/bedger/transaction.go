package bedger

import (
	"bytes"
	"github.com/dgraph-io/badger"
	"github.com/ivanjaros/ijlibs/edger"
	edger_keys "github.com/ivanjaros/ijlibs/edger/keys"
)

type bTxn struct {
	tx     *badger.Txn
	prefix []byte
	kl     int
}

func (t *bTxn) Commit() error {
	return t.tx.Commit()
}

func (t *bTxn) Rollback() error {
	t.tx.Discard()
	return nil
}

func (t *bTxn) BadgerTx() *badger.Txn {
	return t.tx
}

func (t *bTxn) SaveEdges(relType byte, ltrPairs ...[]byte) error {
	keys, err := edger_keys.BuildEdgesForPairs(relType, t.kl, t.prefix, ltrPairs...)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := t.tx.Set(key, nil); err != nil {
			return err
		}
	}

	return nil
}

func (t *bTxn) DeleteEdges(relType byte, ltrPairs ...[]byte) error {
	keys, err := edger_keys.BuildEdgesForPairs(relType, t.kl, t.prefix, ltrPairs...)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := t.tx.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

func (t *bTxn) DeleteItem(relType byte, item []byte) error {
	parents, err := t.LoadParents(relType, item, 1)
	if err != nil {
		return err
	}

	children, err := t.LoadChildren(relType, item, 1)
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

	return t.DeleteEdges(relType, pairs...)
}

func (t *bTxn) LoadParentsUntil(relType byte, child, parent []byte) (edger.Edge, error) {
	return t.loadRange(relType, edger_keys.Ltr, child, parent)
}

func (t *bTxn) LoadChildrenUntil(relType byte, parent, child []byte) (edger.Edge, error) {
	return t.loadRange(relType, edger_keys.Rtl, parent, child)
}

func (t *bTxn) loadRange(relType byte, way byte, left, right []byte) (edger.Edge, error) {
	result := edger.Edge{Item: left}

	itOps := badger.DefaultIteratorOptions
	itOps.PrefetchValues = false
	it := t.tx.NewIterator(itOps)
	defer it.Close()

	rangeReader(it, &result, relType, way, right, t.prefix...)

	return result, nil
}

func (t *bTxn) LoadParents(relType byte, item []byte, maxLevel ...int) (edger.Edge, error) {
	return t.loadRelations(relType, edger_keys.Ltr, item, maxLevel...)
}

func (t *bTxn) LoadChildren(relType byte, item []byte, maxLevel ...int) (edger.Edge, error) {
	return t.loadRelations(relType, edger_keys.Rtl, item, maxLevel...)
}

func (t *bTxn) loadRelations(relType byte, way byte, child []byte, maxLevel ...int) (edger.Edge, error) {
	result := edger.Edge{Item: child}

	max := -1
	if len(maxLevel) > 0 && maxLevel[0] > 0 {
		max = maxLevel[0]
	}

	itOps := badger.DefaultIteratorOptions
	itOps.PrefetchValues = false
	it := t.tx.NewIterator(itOps)
	defer it.Close()

	loopReader(it, &result, relType, way, max, t.prefix...)

	return result, nil
}

func rangeReader(it *badger.Iterator, rel *edger.Edge, relType byte, relOrd byte, until []byte, p ...byte) {
	prefix := append(p, edger_keys.NewPrefix(relType, relOrd, rel.Item)...)
	found := getPrefixed(it, prefix)
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
		rangeReader(it, &rel.Edges[k], relType, relOrd, until, p...)
	}
}

func loopReader(it *badger.Iterator, rel *edger.Edge, relType byte, relOrd byte, maxLoops int, p ...byte) {
	prefix := append(p, edger_keys.NewPrefix(relType, relOrd, rel.Item)...)
	found := getPrefixed(it, prefix)
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
		loopReader(it, &rel.Edges[k], relType, relOrd, maxLoops-1, p...)
	}
}

func getPrefixed(it *badger.Iterator, prefix []byte) [][]byte {
	var matches [][]byte

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		matches = append(matches, it.Item().KeyCopy(nil))
	}

	return matches
}
