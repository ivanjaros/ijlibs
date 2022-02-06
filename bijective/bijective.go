package bijective

import (
	"math"
	"strings"
)

type Converter struct {
	Alphabet string
}

func (c *Converter) Encode(n int64) string {
	if n == 0 {
		return string(c.Alphabet[0])
	}

	var res []string
	base := int64(len(c.Alphabet))

	for n > 0 {
		res = append(res, string(c.Alphabet[n%base]))
		n = int64(math.Floor(float64(n / base)))
	}

	return strRev(strings.Join(res, ""))
}

func (c *Converter) Decode(s string) int64 {
	var i int64
	base := int64(len(c.Alphabet))
	for _, ch := range s {
		var pos int64
		for k, v := range c.Alphabet {
			if v == ch {
				pos = int64(k)
				break
			}
		}
		i = i*base + pos
	}

	return i
}

// https://stackoverflow.com/a/1754209
func strRev(input string) string {
	// Get Unicode code points.
	n := 0
	rn := make([]rune, len(input))
	for _, r := range input {
		rn[n] = r
		n++
	}
	rn = rn[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		rn[i], rn[n-1-i] = rn[n-1-i], rn[i]
	}
	// Convert back to UTF-8.
	return string(rn)
}
