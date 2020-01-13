package runewidth

import "golang.org/x/text/width"

func RuneWidth(ch rune) int {
	switch width.LookupRune(ch).Kind() {
	case width.EastAsianAmbiguous,
		width.EastAsianWide,
		width.EastAsianFullwidth:
		return 2
	case width.Neutral,
		width.EastAsianNarrow,
		width.EastAsianHalfwidth:
		return 1
	default:
		return 0
	}
}
