package gutils

import (
	"strconv"
	"strings"
)

// formVersionKeywords (any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
var formVersionKeywords = map[string]int{
	"dev":   0,
	"alpha": 1,
	"a":     1,
	"beta":  2,
	"b":     2,
	"RC":    3,
	"rc":    3,
	"#":     4,
	"pl":    5,
	"p":     5,
}

var replaceVersionPrefixMap = map[string]string{
	"version": "",
	"Version": "",
	"V":       "",
	"v":       "",
	"-":       ".",
}

// special compare special version forms
func special(form1, form2 string) int {
	found1, found2, len1, len2 := -1, -1, len(form1), len(form2)
	// (any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
	for name, order := range formVersionKeywords {
		if len1 < len(name) {
			continue
		}
		if strings.Compare(form1[:len(name)], name) == 0 {
			found1 = order
			break
		}
	}

	for name, order := range formVersionKeywords {
		if len2 < len(name) {
			continue
		}

		if strings.Compare(form2[:len(name)], name) == 0 {
			found2 = order
			break
		}
	}

	if found1 == found2 {
		return 0
	}

	if found1 > found2 {
		return 1
	}

	return -1
}

// canonicalize
func canonicalize(version string) string {
	ver := []byte(version)
	l := len(ver)
	if l == 0 {
		return ""
	}
	var buf = make([]byte, l*2)
	j := 0
	for i, v := range ver {
		next := uint8(0)
		if i+1 < l { // have the next one
			next = ver[i+1]
		}

		if v == '-' || v == '_' || v == '+' { // replace "-","_","+" to "."
			if j > 0 && buf[j-1] != '.' {
				buf[j] = '.'
				j++
			}
		} else if (next > 0) &&
			(!(next >= '0' && next <= '9') && (v >= '0' && v <= '9')) ||
			(!(v >= '0' && v <= '9') && (next >= '0' && next <= '9')) {
			// insert '.' before and after a non-digit
			buf[j] = v
			j++
			if v != '.' && next != '.' {
				buf[j] = '.'
				j++
			}
			continue
		} else if !((v >= '0' && v <= '9') ||
			(v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z')) { // Non-letters and numbers
			if j > 0 && buf[j-1] != '.' {
				buf[j] = '.'
				j++
			}
		} else {
			buf[j] = v
			j++
		}
	}

	return string(buf[:j])
}

// compare version compare
func compare(origV1, origV2 string) int {
	if origV1 == "" || origV2 == "" {
		if origV1 == "" && origV2 == "" {
			return 0
		}

		if origV1 == "" {
			return -1
		}

		return 1
	}

	ver1, ver2, cmp := "", "", 0
	if origV1[0] == '#' {
		ver1 = origV1
	} else {
		ver1 = canonicalize(origV1)
	}
	if origV2[0] == '#' {
		ver2 = origV2
	} else {
		ver2 = canonicalize(origV2)
	}
	n1, n2 := 0, 0
	for {
		p1, p2 := "", ""
		n1 = strings.IndexByte(ver1, '.')
		if n1 == -1 {
			p1, ver1 = ver1, ""
		} else {
			p1, ver1 = ver1[:n1], ver1[n1+1:]
		}
		n2 = strings.IndexByte(ver2, '.')
		if n2 == -1 {
			p2, ver2 = ver2, ""
		} else {
			p2, ver2 = ver2[:n2], ver2[n2+1:]
		}
		if (p1[0] >= '0' && p1[0] <= '9') && (p2[0] >= '0' && p2[0] <= '9') { // all isdigit
			l1, _ := strconv.Atoi(p1)
			l2, _ := strconv.Atoi(p2)
			if l1 > l2 {
				cmp = 1
			} else if l1 == l2 {
				cmp = 0
			} else {
				cmp = -1
			}
		} else if !(p1[0] >= '0' && p1[0] <= '9') && !(p2[0] >= '0' && p2[0] <= '9') { // all isdigit
			cmp = special(p1, p2)
		} else { // part isdigit
			if p1[0] >= '0' && p1[0] <= '9' { // isdigit
				cmp = special("#N#", p2)
			} else {
				cmp = special(p1, "#N#")
			}
		}
		if cmp != 0 || n1 == -1 || n2 == -1 {
			break
		}
	}

	if cmp == 0 {
		if ver1 != "" {
			if ver1[0] >= '0' && ver1[0] <= '9' {
				cmp = 1
			} else {
				cmp = compare(ver1, "#N#")
			}
		} else if ver2 != "" {
			if ver2[0] >= '0' && ver2[0] <= '9' {
				cmp = -1
			} else {
				cmp = compare("#N#", ver2)
			}
		}
	}

	return cmp
}

// VersionCompare compare the relationship between two version numbers
// The possible operators are: <, lt, <=, le, >, gt, >=, ge, ==, =, eq, !=, <>, ne respectively.
// special version strings these are handled in the following order,
// (any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
// Usage:
// VersionCompare("1.2.3-alpha", "1.2.3RC7", ">=")
// VersionCompare("1.2.3-beta", "1.2.3pl", "lt")
// VersionCompare("1.1_dev", "1.2any", "eq")
func VersionCompare(v1, v2, operator string) bool {
	for k, v := range replaceVersionPrefixMap {
		if strings.Contains(v1, k) {
			v1 = strings.Replace(v1, k, v, -1)
		}

		if strings.Contains(v2, k) {
			v2 = strings.Replace(v2, k, v, -1)
		}
	}

	cmp := compare(v1, v2)
	switch operator {
	case "<", "lt":
		return cmp == -1
	case "<=", "le", "leq":
		return cmp != 1
	case ">", "gt":
		return cmp == 1
	case ">=", "ge", "geq":
		return cmp != -1
	case "==", "=", "eq":
		return cmp == 0
	case "!=", "<>", "ne", "neq":
		return cmp != 0
	default:
		panic("operator: invalid")
	}
}
