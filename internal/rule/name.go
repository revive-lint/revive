// Package rule defines utility functions some revive rules can share.
package rule

import (
	"strings"
	"unicode"
)

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IDS":   true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

// Name returns a different name of struct, var, const, or function if it should be different.
func Name(name string, allowlist, blocklist []string, skipInitialismNameChecks, initialismsAsWords bool) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	ignoreInitWarnings := make(map[string]bool, len(allowlist))
	for _, a := range allowlist {
		ignoreInitWarnings[a] = true
	}

	extraInits := make(map[string]bool, len(blocklist))
	for _, b := range blocklist {
		extraInits[b] = true
	}

	// Split camelCase at any lower->upper transition, at uppercase run transitions, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		switch {
		case i+1 == len(runes):
			eow = true
		case runes[i+1] == '_':
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		case unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]):
			// lower->non-lower
			eow = true
		case i > w && unicode.IsUpper(runes[i]) && i+2 < len(runes) && unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i+2]):
			// upper->upper->lower (e.g. HTTPMethod: HTTP ends at P, Method starts at M)
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		u := strings.ToUpper(word)

		switch {
		case skipInitialismNameChecks:
			if w > 0 && strings.ToLower(word) == word {
				runes[w] = unicode.ToUpper(runes[w])
			}
		case !initialismsAsWords:
			if (commonInitialisms[u] || extraInits[u]) && !ignoreInitWarnings[u] {
				formatStandardInitialism(runes, w, u)
			} else if w > 0 && strings.ToLower(word) == word {
				// already all lowercase, and not the first word, so uppercase the first character.
				runes[w] = unicode.ToUpper(runes[w])
			}
		default:
			applyInitialismsAsWords(runes, w, i, u, extraInits, ignoreInitWarnings)
		}
		w = i
	}
	return string(runes)
}

func formatStandardInitialism(runes []rune, w int, u string) {
	// Keep consistent case, which is lowercase only at the start.
	if w == 0 && unicode.IsLower(runes[w]) {
		u = strings.ToLower(u)
	}
	// Keep lowercase s for IDs
	if u == "IDS" {
		u = "IDs"
	}
	// All the common initialisms are ASCII,
	// so we can replace the bytes exactly.
	copy(runes[w:], []rune(u))
}

func toWordCase(runes []rune, w, i int) {
	if w > 0 {
		runes[w] = unicode.ToUpper(runes[w])
	} else if !unicode.IsUpper(runes[w]) {
		return
	}
	for k := w + 1; k < i; k++ {
		runes[k] = unicode.ToLower(runes[k])
	}
}

func applyInitialismsAsWords(runes []rune, w, i int, u string, extraInits, ignoreInitWarnings map[string]bool) {
	if ignoreInitWarnings[u] {
		return
	}
	if extraInits[u] {
		// Blocklist entries are enforced as standard initialisms
		formatStandardInitialism(runes, w, u)
		return
	}
	toWordCase(runes, w, i)
}
