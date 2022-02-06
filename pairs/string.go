package pairs

// Taken from Gorilla Muxer
func Strings(pairs ...string) (map[string]string, bool) {
	length, err := checkStringPairs(pairs...)
	if err {
		return nil, err
	}
	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m, false
}

func checkStringPairs(pairs ...string) (int, bool) {
	length := len(pairs)
	if length%2 != 0 {
		return length, false
	}
	return length, true
}
