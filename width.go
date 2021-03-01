package readline

import (
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// WidthT means the width type
type WidthT int

var widthCache = map[rune]WidthT{}

// ResetCharWidth resets the cache for the width of characters.
func ResetCharWidth() {
	widthCache = map[rune]WidthT{}
}

// SetCharWidth sets the width of the character into the cache.
func SetCharWidth(c rune, width int) {
	widthCache[c] = WidthT(width)
}

func lenEscaped(c rune) WidthT {
	w := WidthT(3) // '<' + 1-digit + '>'
	for c > 0xF {
		c >>= 4
		w++
	}
	return w
}

func isVariationSelectorLikeRune(ch rune) bool {
	return unicode.Is(unicode.Variation_Selector, ch) ||
		unicode.Is(unicode.Me, ch)
}

func isVariationSelectorLikeStr(s string) bool {
	ch, _ := utf8.DecodeRuneInString(s)
	return isVariationSelectorLikeRune(ch)
}

func isToBeEscaped(ch rune) bool {
	return isVariationSelectorLikeRune(ch) ||
		(ch >= 0x10000 && !SurrogatePairOk) ||
		runewidth.RuneWidth(ch) == 0
}
