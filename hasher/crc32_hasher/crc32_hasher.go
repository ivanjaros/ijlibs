package crc32_hasher

import (
	"github.com/ivanjaros/ijlibs/hasher"
	"hash/crc32"
)

func New() *hasher.Hasher {
	return hasher.New(crc32.NewIEEE())
}
