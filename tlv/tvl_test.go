package tlv

import (
	"encoding/binary"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	var data []byte
	src := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
	values := strings.Split(src, " ")
	ln := make([]byte, 4)
	for _, v := range values {
		binary.BigEndian.PutUint32(ln[:], uint32(len(v)))
		data = append(data, ln...)
		data = append(data, []byte(v)...)
	}

	r := NewReader(data)
	var parsed []string

	for r.Next() {
		parsed = append(parsed, string(r.Read()))
	}

	if result := strings.Join(parsed, " "); result != src {
		t.Fatalf("expected '%s', got '%s'", src, result)
	}
}

func TestWriter(t *testing.T) {
	w := NewWriter(nil)
	src := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
	values := strings.Split(src, " ")
	for _, v := range values {
		w.Write([]byte(v))
	}

	r := NewReader(w.Copy())
	var parsed []string

	for r.Next() {
		parsed = append(parsed, string(r.Read()))
	}

	if result := strings.Join(parsed, " "); result != src {
		t.Fatalf("expected '%s', got '%s'", src, result)
	}
}
