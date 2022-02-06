package tlv

import "encoding/binary"

func NewWriter(data []byte) *Writer {
	return &Writer{data: data}
}

type Writer struct {
	data []byte
}

func (w *Writer) Write(data []byte) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size[:], uint32(len(data)))
	w.data = append(w.data, size...)
	w.data = append(w.data, data...)
}

func (w *Writer) Result() []byte {
	return w.data
}

func (w *Writer) Copy() []byte {
	cp := make([]byte, len(w.data))
	copy(cp, w.data)
	return cp
}

func (w *Writer) Reset() {
	w.data = nil
}
