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

func (c _PackedColorCode) WriteTo(w io.Writer) (int64, error) {
	if c < 0 {
		return 0, nil
	}
	ofs := "\x1B["
	n := int64(0)
	for ; c > 0; c >>= colorCodeBitSize {
		_n, err := fmt.Fprintf(w, "%s%d", ofs, c&colorCodeMask)
		n += int64(_n)
		if err != nil {
			return n, err
		}
		ofs = ";"
	}
	_n, err := w.Write([]byte{'m'})
	n += int64(_n)
	return n, err
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
			color.WriteTo(B.Out)
		}
		ch.Moji.PrintTo(B.Out)
	}
	if color != defaultColor {
		defaultColor.WriteTo(B.Out)
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Moji.Width()
	}
	return
}
