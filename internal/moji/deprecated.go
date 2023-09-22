package moji

// Deprecated: GetCharWidth returns the width of the character.
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
			width = WidthT(runeWidth(n))
			if width == 0 {
				width = lenEscaped(n)
			}
		}
		widthCache[n] = width
	}
	return width
}
