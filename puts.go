package readline

import (
	"fmt"
	"io"
)

func (B *Buffer) backspace(n WidthT) {
	if n > 1 {
		fmt.Fprintf(B.Out, "\x1B[%dD", n)
	} else if n == 1 {
		B.Out.WriteByte('\b')
	}
}

func (B *Buffer) eraseline() {
	io.WriteString(B.Out, "\x1B[0K")
}

type _Range []_Cell

const (
	colorCodeBitSize = 8
	colorCodeMask    = (1<<colorCodeBitSize - 1)
)

func SGR1(n1 int) int     { return n1 }
func SGR2(n1, n2 int) int { return n1 | (n2 << colorCodeBitSize) }

func SGR3(n1, n2, n3 int) int {
	return n1 |
		(n2 << colorCodeBitSize) |
		(n3 << (colorCodeBitSize * 2))
}

func SGR4(n1, n2, n3, n4 int) int {
	return n1 |
		(n2 << colorCodeBitSize) |
		(n3 << (colorCodeBitSize * 2)) |
		(n4 << (colorCodeBitSize * 3))
}

func putColor(w io.Writer, c _PackedColorCode) {
	if c < 0 {
		return
	}
	ofs := "\x1B["
	for ; c > 0; c >>= colorCodeBitSize {
		fmt.Fprintf(w, "%s%d", ofs, c&colorCodeMask)
		ofs = ";"
	}
	w.Write([]byte{'m'})
}

func (B *Buffer) Write(b []byte) (int, error) {
	return B.Out.Write(b)
}

func (B *Buffer) puts(s []_Cell) _Range {
	defaultColor := _PackedColorCode(B.RefreshColor())
	color := defaultColor
	for _, ch := range s {
		if ch.color != color {
			color = ch.color
			putColor(B.Out, color)
		}
		ch.Moji.PrintTo(B.Out)
	}
	if color != defaultColor {
		putColor(B.Out, defaultColor)
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Moji.Width()
	}
	return
}
