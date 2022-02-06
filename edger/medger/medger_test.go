//  Copyright (c) Ivan Jaroš. All rights reserved.
//  Any action regarding this code(manipulation, alteration, transformation...) is hereby prohibited
//  without previous explicit consent of the author.
//  This excludes any imported third-party libraries whose usage is guided by their respective licenses.

package medger

import (
	"bytes"
	"testing"
)

func TestMedger(t *testing.T) {
	m, _ := New(3)
	rt := byte(0)

	pairs := [][]byte{
		[]byte("001"), []byte("010"), []byte("010"), []byte("100"),
		[]byte("002"), []byte("020"), []byte("020"), []byte("200"),
		[]byte("003"), []byte("030"), []byte("030"), []byte("300"),
	}

	if err := m.SaveEdges(rt, pairs...); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		child := i
		if child > 0 {
			child = i * 4
		}

		edges, err := m.LoadParents(rt, pairs[child])
		if err != nil {
			t.Fatal(err)
		}

		ids := edges.GetIds()
		if len(ids) != 3 {
			t.Fatalf("expected 3 edges, got %d", len(ids))
		}

		for k, id := range ids {
			if k == 2 {
				k++
			}
			if bytes.Equal(id, pairs[child+k]) == false {
				t.Fatalf("parent '%s' does not match '%s'", string(id), string(pairs[child+k]))
			}
		}
	}

	if err := m.SaveEdges(rt, []byte("100"), []byte("200"), []byte("200"), []byte("300")); err != nil {
		t.Fatal(err)
	}

	edges, err := m.LoadParents(rt, pairs[0])
	if err != nil {
		t.Fatal(err)
	}

	ids := edges.GetIds()
	if len(ids) != 5 {
		t.Fatalf("expected 5 edges, got %d", len(ids))
	}

	if bytes.Equal(ids[3], pairs[7]) == false {
		t.Fatalf("parent '%s' does not match '%s'", string(ids[3]), string(pairs[7]))
	}

	if bytes.Equal(ids[4], pairs[11]) == false {
		t.Fatalf("parent '%s' does not match '%s'", string(ids[4]), string(pairs[11]))
	}
}
