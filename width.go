package readline

import (
	"os"

	"github.com/mattn/go-runewidth"
)

type width_t int

var widthCache = map[rune]width_t{}

func ResetCharWidth() {
	widthCache = map[rune]width_t{}
}

func SetCharWidth(c rune, width int) {
	widthCache[c] = width_t(width)
}

func lenEscaped(c rune) width_t {
	w := width_t(3) // '<' + 1-digit + '>'
	for c > 0xF {
		c >>= 4
		w++
	}
	return w
}

var isWindowsTerminal = (os.Getenv("WT_SESSION") != "" &&
	os.Getenv("WT_PROFILE_ID") != "")

func GetCharWidth(n rune) width_t {
	if n < ' ' {
		return 2 // ^X
	}
	if n == 0x7F {
		return 4 // <7F>
	}
	width, ok := widthCache[n]
	if !ok {
		if n > 0x10000 && !isWindowsTerminal {
			width = lenEscaped(n)
		} else {
			if isWindowsTerminal && runewidth.IsAmbiguousWidth(n) {
				width = 1
			} else {
				width = width_t(runewidth.RuneWidth(n))
				if width == 0 {
					width = lenEscaped(n)
				}
			}
		}
		widthCache[n] = width
	}
	return width
}

func GetStringWidth(s string) width_t {
	width := width_t(0)
	for _, ch := range s {
		width += GetCharWidth(ch)
	}
	return width
}
