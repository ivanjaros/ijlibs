package bedger

import (
	"emperror.dev/errors"
	"github.com/dgraph-io/badger"
	"github.com/ivanjaros/ijlibs/edger"
)

// allows hooking up into the underlying Badger transaction
type BadgerTransaction interface {
	edger.Transaction
	BadgerTx() *badger.Txn
}

// allows hooking up into existing Badger instance with
// Badger transaction being handled outside of the edger.
type BadgerEdger interface {
	edger.Edger
	BadgerTx(*badger.Txn) BadgerTransaction
}

// Creates new edger instance backed up by Badger database.
// Prefix can be optionally added into each key to prevent key collisions in case
// the badger instance is being used elsewhere.
func New(db *badger.DB, keySize int, prefix ...byte) (*bEdger, error) {
	if db == nil {
		return nil, errors.New("no badger connection provided")
	}
	if keySize < 1 {
		return nil, errors.New("invalid key length")
	}
	return &bEdger{b: db, prefix: prefix}, nil
}

type bEdger struct {
	b      *badger.DB
	prefix []byte
	kln    int
}

func (e *bEdger) Transaction() edger.Transaction {
	return e.newTx(false)
}

func (e *bEdger) newTx(ro bool) edger.Transaction {
	return &bTxn{
		tx:     e.b.NewTransaction(!ro),
		prefix: e.prefix,
	}
}

func (e *bEdger) BadgerTx(tx *badger.Txn) BadgerTransaction {
	return &bTxn{
		tx:     tx,
		prefix: e.prefix,
	}
}

func (e *bEdger) SaveEdges(relType byte, ltrPairs ...[]byte) error {
	tx := e.Transaction()
	defer tx.Rollback()
	if err := tx.SaveEdges(relType, ltrPairs...); err != nil {
		return err
	}
	return tx.Commit()
}

func (e *bEdger) DeleteEdges(relType byte, ltrPairs ...[]byte) error {
	tx := e.Transaction()
	defer tx.Rollback()
	if err := tx.DeleteEdges(relType, ltrPairs...); err != nil {
		return err
	}
	return tx.Commit()
}

func (e *bEdger) DeleteItem(relType byte, item []byte) error {
	tx := e.Transaction()
	defer tx.Rollback()
	if err := tx.DeleteItem(relType, item); err != nil {
		return err
	}
	return tx.Commit()
}

func (e *bEdger) LoadParents(relType byte, item []byte, maxLevel ...int) (edger.Edge, error) {
	tx := e.newTx(true)
	defer tx.Rollback()
	return tx.LoadParents(relType, item, maxLevel...)
}

func (e *bEdger) LoadChildren(relType byte, item []byte, maxLevel ...int) (edger.Edge, error) {
	tx := e.newTx(true)
	defer tx.Rollback()
	return tx.LoadChildren(relType, item, maxLevel...)
}

func (e *bEdger) LoadParentsUntil(relType byte, child, parent []byte) (edger.Edge, error) {
	tx := e.newTx(true)
	defer tx.Rollback()
	return tx.LoadParentsUntil(relType, child, parent)
}

func (e *bEdger) LoadChildrenUntil(relType byte, parent, child []byte) (edger.Edge, error) {
	tx := e.newTx(true)
	defer tx.Rollback()
	return tx.LoadChildrenUntil(relType, parent, child)
}
