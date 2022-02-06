package edger_keys

import (
	"errors"
)

const (
	Ltr = byte(0)
	Rtl = byte(1)
)

var ErrOddPairs = errors.New("odd number of values provided")

// counts the number of key-value pairs and returns bool if all keys have a value.
func CountPairs(pairs ...[]byte) (int, bool) {
	length := len(pairs)
	if length%2 != 0 {
		return length, false
	}
	return length / 2, true
}

// applies function on each pair until there are no more pairs or the function returned error
func LoopPairs(f func(k, v []byte) error, pairs ...[]byte) error {
	if _, ok := CountPairs(pairs...); ok == false {
		return ErrOddPairs
	}
	for i := 0; i < len(pairs); i += 2 {
		if err := f(pairs[i], pairs[i+1]); err != nil {
			return err
		}
	}
	return nil
}

// left and right are expected to be the same length, this function does no validation on its own
func NewKey(relType, relOrd byte, left, right []byte) []byte {
	return append([]byte{relType, relOrd}, append(left, right...)...)
}

func NewPrefix(relType, relOrd byte, item []byte) []byte {
	return append([]byte{relType, relOrd}, item...)
}

func NewEdge(relType byte, left, right []byte) ([]byte, []byte) {
	return NewKey(relType, Ltr, left, right), NewKey(relType, Rtl, right, left)
}

// this function expects the key to have even number of bytes and be at least 6 bytes long
// it does no validation on its own
func GetEdges(k []byte) (left, right []byte) {
	ln := (len(k) - 2) / 2
	left = k[2:ln]
	right = k[ln+2:]
	return
}

func BuildEdgesForPairs(relType byte, keyLength int, prefix []byte, pairs ...[]byte) ([][]byte, error) {
	if _, ok := CountPairs(pairs...); ok == false {
		return nil, ErrOddPairs
	}

	keys := make([][]byte, 0, len(pairs))

	err := LoopPairs(func(k, v []byte) error {
		if len(k) != keyLength || len(v) != keyLength {
			return errors.New("invalid key length")
		}
		l, r := NewEdge(relType, k, v)
		if len(prefix) > 0 {
			l = append(prefix, l...)
			r = append(prefix, r...)
		}
		keys = append(keys, l, r)
		return nil
	}, pairs...)

	return keys, err
}
