package reader

import (
	"unicode"
	"unicode/utf8"
	"strings"
)
	
func trimLeftSpace(s string) string {
	// stolen from the Go standard library (first half of TrimSpace)
	// Fast path for ASCII: look for the first ASCII non-space byte
	asciiSpace := [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return strings.TrimLeftFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}
	return s[start:]
}

func trimLeftSpaceComment(s string) string {
	ss := trimLeftSpace(s)
	for len(ss) > 1 && ss[:2] == ";;" {
		// while we got leading comments
		// remove leading comment
		nl := strings.Index(ss, "\n")
		if nl < 0 {
			// no NL, so really we have nothing - bail
			return ""
		}
		// bump the content to the NL
		ss = trimLeftSpace(ss[nl:])
	}
	return ss
}	
