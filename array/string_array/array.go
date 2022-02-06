package string_array

// Some logic has been taken from https://github.com/golang/go/wiki/SliceTricks

import (
	"math/rand"
	"sort"
)

func Contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// filters out duplicates and returns sorted slice
func Unique(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}

	sort.Strings(arr)
	j := 0
	for i := 1; i < len(arr); i++ {
		if arr[j] == arr[i] {
			continue
		}
		j++
		arr[j] = arr[i]
	}

	return arr[:j+1]
}

// removes duplicates from arr without changing order
// ie. first value will be kept in its place, duplicates will be removed.
func UniqueUnsorted(arr []string) []string {
	prod := make([]string, 0, len(arr))

	for k := range arr {
		if Contains(arr[k], prod) == false {
			prod = append(prod, arr[k])
		}
	}

	return prod
}

// indicates if array contains unique values
func IsUnique(arr []string) bool {
	return len(arr) == len(Unique(arr))
}

// returns all values from a, that ARE present in b
func Intersect(a, b []string) []string {
	c := make([]string, 0, len(a))
	for _, v := range a {
		if Contains(v, b) {
			c = append(c, v)
		}
	}
	return Unique(c)
}

// returns all values from a that ARE NOT present in b
func Diff(a, b []string) []string {
	c := make([]string, 0, len(a))
	for _, v := range a {
		if Contains(v, b) == false {
			c = append(c, v)
		}
	}
	return Unique(c)
}

func Copy(a []string) []string {
	b := make([]string, len(a))
	copy(b, a)
	return b
}

// compares two slices and determines if all values are present in both of them.
// this method sorts values so order is ignored.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCp := Copy(a)
	bCp := Copy(b)
	sort.Strings(aCp)
	sort.Strings(bCp)
	return EqualStrict(aCp, bCp)
}

// compares two slices and determines if all values are present in both of them.
// this method takes order into account.
func EqualStrict(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if a[k] != b[k] {
			return false
		}
	}

	return true
}

// filters out values from a and returns changed a
func Filter(a []string, keep func(value string) bool) []string {
	if len(a) == 0 {
		return a
	}

	b := make([]string, 0, len(a))
	for _, v := range a {
		if keep(v) {
			b = append(b, v)
		}
	}

	return b
}

// removes the last item in a and returns changed a and removed value
func Pop(a []string) ([]string, string) {
	if len(a) == 0 {
		return a, ""
	}
	var rem string
	rem, a = a[len(a)-1], a[:len(a)-1]
	return a, rem
}

// removes first item from a and returns changed a and removed value
func Shift(a []string) ([]string, string) {
	if len(a) == 0 {
		return a, ""
	}
	var rem string
	rem, a = a[0], a[1:]
	return a, rem
}

// adds value at the beginning of a and returns changed a
func UnShift(v string, a []string) []string {
	if len(a) == 0 {
		return a
	}
	a = append([]string{v}, a...)
	return a
}

// reverse values in a
func Reverse(a []string) {
	if len(a) > 1 {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	}
}

// randomizes order of values in a
func Shuffle(a []string) {
	if len(a) > 1 {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	}
}

// splits a into multiple smaller chunks with each chunk having up to size items
func Chunk(a []string, size int) [][]string {
	var chunked [][]string

	for size < len(a) {
		a, chunked = a[size:], append(chunked, a[0:size:size])
	}

	return append(chunked, a)
}

func Join(a ...[]string) []string {
	if len(a) == 0 {
		return nil
	}
	var c int
	for k := range a {
		c += len(a[k])
	}
	out := make([]string, 0, c)
	for k := range a {
		out = append(out, a[k]...)
	}
	return out
}

func Delete(src []string, values ...string) []string {
	return Diff(src, values)
}
