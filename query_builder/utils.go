package query_builder

import "strings"

func StringPairs(pairs ...string) (map[string]string, bool) {
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

// Helper function for human/visual testing.
// test https://play.golang.org/p/lHkNUwhLBzO
func ConditionToSqlString(c ICondition, enclose bool) string {
	expr := c.Column()
	if tbl := c.Table(); tbl != "" {
		expr = tbl + TableColumnDivider + expr
	}

	out := append([]string{}, expr, string(c.Operator()), "?")

	if c.IsNegated() {
		out = append([]string{"NOT"}, out...)
	}

	var wrap bool

	for _, v := range c.Siblings() {
		out = append(out, "AND", ConditionToSqlString(v, true))
	}

	if alts := c.Alternatives(); len(alts) > 0 {
		wrap = true
		out = append([]string{"("}, out...)
		out = append(out, ")")
		for _, v := range alts {
			out = append(out, "OR", "(", ConditionToSqlString(v, true), ")")
		}
	}

	joined := strings.Join(out, " ")

	if wrap && enclose {
		joined = "( " + joined + " )"
	}

	return joined
}
