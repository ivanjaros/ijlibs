package sha1_hasher

import (
	"crypto/sha1"
	"github.com/ivanjaros/ijlibs/hasher"
)

func New() *hasher.Hasher {
	return hasher.New(sha1.New())
}
