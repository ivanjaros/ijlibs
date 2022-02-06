package hasher

import (
	"encoding/hex"
	"hash"
)

type Hasher struct {
	alg hash.Hash
}

func New(alg hash.Hash) *Hasher {
	return &Hasher{alg: alg}
}

func (h *Hasher) String(in string) string {
	return h.Bytes([]byte(in))
}

func (h *Hasher) Bytes(in []byte) string {
	h.alg.Reset()
	h.alg.Write(in)
	return hex.EncodeToString(h.alg.Sum(nil))
}

func HashString(alg hash.Hash, value string) string {
	return HashBytes(alg, []byte(value))
}

func HashBytes(alg hash.Hash, value []byte) string {
	alg.Reset()
	alg.Write(value)
	return hex.EncodeToString(alg.Sum(nil))
}
