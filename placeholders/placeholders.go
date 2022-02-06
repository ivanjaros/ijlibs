package placeholders

import (
	"strings"
)

func Columns(count int) string {
	if count < 1 {
		return ""
	}

	var result string
	for i := 0; i < count; i++ {
		result += ",?"
	}
	return result[1:]
}

func Groups(cols int, groups int) string {
	if cols < 1 {
		return "()"
	}

	var result string
	for i := 0; i < groups; i++ {
		result += ",("
		result += Columns(cols)
		result += ")"
	}
	return result[1:]
}

func Group(count int) string {
	return Groups(count, 1)
}

func StringArguments(ids []string) []interface{} {
	args := make([]interface{}, 0, len(ids))
	for k := range ids {
		args = append(args, ids[k])
	}
	return args
}

// takes up multiple strings. the one that is "%" or "_" will be escaped.
// returns the string as joined values with escaped individual values.
func Like(values ...string) string {
	for k, v := range values {
		if v == "%" || v == "_" {
			continue
		}
		values[k] = strings.ReplaceAll(v, "%", "\\%")
		values[k] = strings.ReplaceAll(v, "_", "\\_")
	}
	return strings.Join(values, "")
}
