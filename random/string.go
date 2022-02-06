package random

import (
	"math/rand"
	"time"
)

type CharSet string

const (
	NumSet         CharSet = `0123456789`
	NumSetSafe     CharSet = `23456789`
	AlphaSetLC     CharSet = `abcdefghijklmnopqrstuvwxyz`
	AlphaSetSafeLC CharSet = `abcdefghjkmnpqrstuvwxyz`
	AlphaSetUC     CharSet = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	AlphaSetSafeUC CharSet = `ABCDEFGHJKLMNPQRSTUVWXYZ`
	// https://www.owasp.org/index.php/Password_special_characters
	OwaspSpecial CharSet = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func String(ln int, charSet0 CharSet, charSetN ...CharSet) string {
	chars := []rune(charSet0)
	for k := range charSetN {
		chars = append(chars, []rune(charSetN[k])...)
	}
	l := len(chars)
	randStr := make([]rune, ln)
	for k := range randStr {
		randStr[k] = chars[rnd.Intn(l)]
	}
	return string(randStr)
}

func SafeString(len int) string {
	return String(len, NumSetSafe, AlphaSetSafeLC, AlphaSetSafeUC)
}

func StandardString(len int) string {
	return String(len, NumSet, AlphaSetLC, AlphaSetUC)
}

func StringRange(min, max int) string {
	return SafeString(IntegerRange(min, max))
}
