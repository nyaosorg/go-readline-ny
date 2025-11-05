package moji

import (
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// WidthT means the width type
type WidthT int

func lenEscaped(c rune) WidthT {
	w := WidthT(3) // '<' + 1-digit + '>'
	for c > 0xF {
		c >>= 4
		w++
	}
	return w
}

func isVariationSelectorLike(ch rune) bool {
	if !VariationSequenceOk {
		return false
	}
	return unicode.Is(unicode.Variation_Selector, ch) ||
		unicode.Is(unicode.Me, ch)
}

func AreVariationSelectorLike(s string) bool {
	if !VariationSequenceOk {
		return false
	}
	ch, _ := utf8.DecodeRuneInString(s)
	return isVariationSelectorLike(ch)
}

func isToBeEscaped(ch rune) bool {
	return isVariationSelectorLike(ch) ||
		(ch >= 0x10000 && !SurrogatePairOk) ||
		runewidth.RuneWidth(ch) == 0
}
