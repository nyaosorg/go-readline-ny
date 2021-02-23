package readline

import (
	"github.com/mattn/go-runewidth"
	"unicode"
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

func isVariationSelector(ch rune) bool {
	return unicode.Is(unicode.Variation_Selector, ch)
}

func isToBeEscaped(ch rune) bool {
	return isVariationSelector(ch) ||
		(ch >= 0x10000 && !SurrogatePairOk) ||
		runewidth.RuneWidth(ch) == 0
}

// GetCharWidth returns the width of the character.
func GetCharWidth(n rune) WidthT {
	if n < ' ' {
		return 2 // ^X
	}
	if n == 0x7F {
		return 4 // <7F>
	}
	width, ok := widthCache[n]
	if !ok {
		if isToBeEscaped(n) {
			width = lenEscaped(n)
		} else {
			if TreatAmbiguousWidthAsNarrow && runewidth.IsAmbiguousWidth(n) {
				width = 1
			} else {
				width = WidthT(runewidth.RuneWidth(n))
				if width == 0 {
					width = lenEscaped(n)
				}
			}
		}
		widthCache[n] = width
	}
	return width
}

// GetStringWidth returns the width of the string.
func GetStringWidth(s string) WidthT {
	width := WidthT(0)
	for _, ch := range s {
		width += GetCharWidth(ch)
	}
	return width
}
