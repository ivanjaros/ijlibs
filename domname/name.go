package domname

import (
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"strings"
)

// identifies domain information.
// root is the alpha tld like "com", "uk", "net", "org" and so on.
// suffix is special case for domains like "co.uk", "co.za" and so on, but it might be empty.
// rest is the rest of the domain that prefixes the root and/or the suffix if applicable.
//
// for example:
//  "com" will return "com", "", ""
//  ".com" will return "com", "", ""
//  "foo.com" will return "com", "", "foo"
//  "foo.co.uk" will return "uk", "co.uk", "foo"
//  "ns1.foo.com" will return "com", "", "ns1.foo"
//  "ns1.foo.co.uk" will return "uk", "co.uk", "ns1.foo"
//
// note that "co.uk" will return "", "", "co.uk" but ".co.uk" will return "uk", "co.uk", ""
// for these cases where "complex" tld is used, and it is know to be TLD or expected it to be,
// use the IdentifyTLD() instead.
//
// if root or suffix are not identified, they will be empty and rest will contain the entire domain.
// this function works with lower-case domain names in their idn form.
//
func Identify(domain string) (root string, suffix string, rest string) {
	used := domain
	var prefixed bool

	if strings.Contains(domain, ".") == false {
		prefixed = true
		used = "foobarbaznone." + domain
	}

	if len(used) > 1 && used[0] == '.' {
		prefixed = true
		used = "foobarbaznone" + domain
	}

	// we have to use this approach with no default find rule,
	// otherwise non-existing tlds will return valid results via widlcard rule
	// and when we'll try to process the tld with information from the returned rule
	// we'll get index out of range panic since the domain does not actually match the widlcard rule.
	info, err := publicsuffix.ParseFromListWithOptions(publicsuffix.DefaultList, used, &publicsuffix.FindOptions{IgnorePrivate: false})
	if err != nil {
		return "", "", domain
	}

	root = info.TLD

	rest = info.String()
	rest = rest[:len(rest)-len(root)-1] // -1 for the dot offset

	if info.Rule.Length == 2 {
		suffix = info.TLD
		root = strings.Split(root, ".")[1]
	}

	if prefixed {
		rest = ""
	}

	return
}

// identifies TLD, same as Identify() but only the TLD is expected
func IdentifyTLD(tld string) (root string, suffix string) {
	root, suffix, _ = Identify("." + tld)
	return
}
