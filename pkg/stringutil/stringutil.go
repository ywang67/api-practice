// package stringutil contains helper methods for strings that aren't present in the stdlib strings package.
package stringutil

import (
	"strings"
	"unicode"
)

func InterfaceSliceContains(s []interface{}, e string) bool {
	for _, a := range s {
		if a.(string) == e {
			return true
		}
	}
	return false
}

// Difference returns the elements in 'a' that aren't in 'b'.
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// After: get the sub string after a string
func After(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

// SliceContains returns if a slice of strings contains the given string
func SliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SliceContains returns if a slice of strings contains the given string
func SliceIndex(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

// SliceSetEquals returns true if the sets of two slices of strings are equivalent (same elements, any order)
func SliceSetEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	// for each string key, we count the copies in a and copies in b. if they're the same, the slices are equivalent.
	// but rather than making two maps to count, we can simply _add_ for `a` and _subtract_ for `b`.
	counts := make(map[string]int, 2*len(a)) // worst case scenario: A and B have no shared keys.
	for _, k := range a {
		counts[k]++
	}
	for _, k := range b {
		counts[k]--
	}

	// c < 0: too many copies in b
	// c > 0: too many copies in a
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}
	return true
}

// ContainsByte checks for the presence of b in s. Like strings.ContainsRune,
// but for bytes.
func ContainsByte(s string, b byte) bool { return strings.IndexByte(s, b) != -1 }

// IsASCIIAlphaNumeric tests whether a string is entirely alphanumeric ASCII characters.
// Equivalent to matching the regexp [a-zA-Z0-9].
func IsASCIIAlphaNumeric(s string) bool {
	for i := range s {
		b := s[i]
		switch {
		case '0' <= b && b <= '9', 'a' <= b && b <= 'z', 'A' <= b && b <= 'Z':
			continue
		default:
			return false
		}
	}
	return true
}

// Unpack destructures a slice of strings into separate assigned variable references
func Unpack(s []string, vars ...*string) {
	for i, str := range s {
		if i >= len(vars) {
			return
		}
		*vars[i] = str
	}
}

// TrimSpace removes all outer whitespaces from strings
func TrimSpace(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

// TrimSpecialCharter removes all outer special charter from strings
func TrimSpecialCharter(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func HyphenToCamelCase(str string) string {
	var tempStr string
	upper := false
	for _, c := range str {
		char := string(c)
		if upper {
			tempStr += strings.ToUpper(char)
			upper = false
		} else {
			if char == "-" {
				upper = true
			} else {
				tempStr += char
			}
		}
	}
	return tempStr
}
