package readline

import (
	"github.com/mattn/go-runewidth"
)

// WidthT means the width type
type WidthT int

var widthCache = map[rune]WidthT{}

func ResetCharWidth() {
	widthCache = map[rune]WidthT{}
}

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

func GetCharWidth(n rune) WidthT {
	if n < ' ' {
		return 2 // ^X
	}
	if n == 0x7F {
		return 4 // <7F>
	}
	width, ok := widthCache[n]
	if !ok {
		if n > 0x10000 && !SurrogatePairOk {
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

func GetStringWidth(s string) WidthT {
	width := WidthT(0)
	for _, ch := range s {
		width += GetCharWidth(ch)
	}
	return width
}
