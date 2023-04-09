package readline

// GetCharWidth returns the width of the character. [Deprecated]
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
			width = WidthT(wtRuneWidth.RuneWidth(n))
			if width == 0 {
				width = lenEscaped(n)
			}
		}
		widthCache[n] = width
	}
	return width
}

// GetStringWidth returns the width of the string. [Deprecated]
func GetStringWidth(s string) WidthT {
	width := WidthT(0)
	for _, m := range StringToMoji(s) {
		width += m.Width()
	}
	return width
}
