package moji

import (
	"unicode"
	"unicode/utf8"

	"github.com/nyaosorg/go-readline-ny/internal/termcheck"
	// rw "github.com/nyaosorg/go-readline-ny/internal/orgxwidth"
	rw "github.com/nyaosorg/go-readline-ny/internal/runewidth"
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

var condition interface {
	RuneWidth(rune) int
} = rw.New(termcheck.AmbiguousIsWide)

func getWidth(r rune) WidthT {
	w, ok := widthCache[r]
	if !ok {
		w = WidthT(condition.RuneWidth(r))
		widthCache[r] = w
	}
	return w
}

func lenEscaped(c rune) WidthT {
	w := WidthT(3) // '<' + 1-digit + '>'
	for c > 0xF {
		c >>= 4
		w++
	}
	return w
}

func isVariationSelectorLike(ch rune) bool {
	if !termcheck.VariationSequenceOk {
		return false
	}
	return unicode.Is(unicode.Variation_Selector, ch) ||
		unicode.Is(unicode.Me, ch)
}

func AreVariationSelectorLike(s string) bool {
	if !termcheck.VariationSequenceOk {
		return false
	}
	ch, _ := utf8.DecodeRuneInString(s)
	return isVariationSelectorLike(ch)
}

func isToBeEscaped(ch rune) bool {
	return isVariationSelectorLike(ch) ||
		(ch >= 0x10000 && !termcheck.SurrogatePairOk) ||
		getWidth(ch) == 0
}
