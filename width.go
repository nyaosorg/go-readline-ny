package readline

import (
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

var x = struct{}{}

var variationSelector = map[rune]struct{}{
	// FVS: Mongolian Free Variation Selector
	'\u180B': x,
	'\u180C': x,
	'\u180D': x,
	// SVS: Standardized Variation Sequence
	'\uFE00': x,
	'\uFE01': x,
	'\uFE02': x,
	'\uFE03': x,
	'\uFE04': x,
	'\uFE05': x,
	'\uFE06': x,
	'\uFE07': x,
	'\uFE08': x,
	'\uFE09': x,
	'\uFE0A': x,
	'\uFE0B': x,
	'\uFE0C': x,
	'\uFE0D': x,
	'\uFE0E': x,
	'\uFE0F': x,
	// IVS: Ideographic Variation Sequence
	'\U000E0100': x,
	'\U000E0101': x,
	'\U000E0102': x,
	'\U000E0103': x,
	'\U000E0104': x,
	'\U000E0105': x,
	'\U000E0106': x,
	'\U000E0107': x,
	'\U000E0108': x,
	'\U000E0109': x,
	'\U000E010A': x,
	'\U000E010B': x,
	'\U000E010C': x,
	'\U000E010D': x,
	'\U000E010E': x,
	'\U000E01EF': x,
}

func isVariationSelector(ch rune) bool {
	_, ok := variationSelector[ch]
	return ok
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
