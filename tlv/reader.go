package tlv

import "encoding/binary"

func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

type Reader struct {
	data    []byte
	pointer int
	inited  bool
}

func (r *Reader) Read() []byte {
	if r.inited == false {
		r.inited = true
	}

	if r.pointer+4 >= len(r.data) {
		return nil
	}

	size := int(binary.BigEndian.Uint32(r.data[r.pointer : r.pointer+4]))
	if (r.pointer + 4 + size) > len(r.data) {
		return nil
	}

	return r.data[r.pointer+4 : r.pointer+4+size]
}

func (r *Reader) Next() bool {
	if r.inited == false {
		r.inited = true
		return true
	}

	if r.pointer+4 >= len(r.data) {
		return false
	}

	size := int(binary.BigEndian.Uint32(r.data[r.pointer : r.pointer+4]))
	r.pointer += 4
	r.pointer += size

	return r.pointer < len(r.data)
}

func (r *Reader) Reset() {
	r.pointer = 0
	r.inited = false
}
