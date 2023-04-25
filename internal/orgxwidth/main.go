package orgxwidth

import (
	"unicode"

	"golang.org/x/text/width"
)

type Condition struct {
	AmbiguousIsWide bool
}

func New(ambiguousIsWide bool) Condition {
	return Condition{ AmbiguousIsWide:ambiguousIsWide }
}

func (C Condition) RuneWidth(r rune) int {
	if !unicode.IsPrint(r) {
		return 0
	}
	switch width.LookupRune(r).Kind() {
	case width.Neutral, width.EastAsianNarrow, width.EastAsianHalfwidth:
		return 1
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	case width.EastAsianAmbiguous:
		if C.AmbiguousIsWide {
			return 2
		} else {
			return 1
		}
	default:
		return 0
	}
}
