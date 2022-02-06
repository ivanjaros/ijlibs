package str

import (
	"unicode/utf8"
)

// https://stackoverflow.com/a/70793821/17644273
func Cut(name string, ln int) string {
	if utf8.RuneCountInString(name) > ln {
		// convert string to rune slice
		chars := []rune(name)
		// index the rune slice and convert back to string
		name = string(chars[:ln])
	}
	return name
}
