package fanner

import "testing"

func TestFanIn(t *testing.T) {
	fin := NewIn(2)

	l := fin.Listen()
	a := make(chan interface{})
	b := make(chan interface{})

	fin.Add(a)
	fin.Add(b)

	a <- "a"
	b <- "b"

	fin.Close()

	if v := <-l; v != "a" {
		t.Fatalf("expected 'a', got '%v'", v)
	}
	if v := <-l; v != "b" {
		t.Fatalf("expected 'b', got '%v'", v)
	}
}

func TestPipeIn(t *testing.T) {
	fin1 := NewIn(2)
	a1 := make(chan interface{})
	b1 := make(chan interface{})
	fin1.Add(a1)
	fin1.Add(b1)
	l1 := fin1.Listen()

	fin2 := NewIn(2)
	a2 := make(chan interface{})
	b2 := make(chan interface{})
	fin2.Add(a2)
	fin2.Add(b2)
	l2 := fin2.Listen()

	piped := NewIn(4)
	piped.Add(l1)
	piped.Add(l2)
	pipe := piped.Listen()

	a1 <- "a1"
	b1 <- "b1"
	a2 <- "a2"
	b2 <- "b2"

	fin1.Close()
	fin2.Close()
	piped.Close()

	expect := []string{"a1", "b1", "a2", "b2"}

	for {
		select {
		case v, ok := <-pipe:
			if ok == false {
				t.Fatal("pipe is closed")
			}

			var found bool
			for k := range expect {
				if expect[k] == v {
					expect = append(expect[:k], expect[k+1:]...)
					found = true
					break
				}
			}
			if found == false {
				t.Fatalf("unexpected value '%v'", v)
			}

			if len(expect) == 0 {
				return
			}

		default:
			t.Fatal("pipe is empty")
		}
	}
}
