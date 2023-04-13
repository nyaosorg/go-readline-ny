package readline

import (
	"fmt"
	"io"
)

func (B *Buffer) eraseline() {
	io.WriteString(B.Out, "\x1B[0K")
}

type _Range []Cell

const (
	colorCodeBitSize = 8
	colorCodeMask    = (1<<colorCodeBitSize - 1)
)

type ColorSequence int64

func SGR1(n1 int) ColorSequence {
	return ColorSequence(1) |
		(ColorSequence(n1) << colorCodeBitSize)
}

func SGR2(n1, n2 int) ColorSequence {
	return ColorSequence(2) |
		(ColorSequence(n1) << colorCodeBitSize) |
		(ColorSequence(n2) << (colorCodeBitSize * 2))
}

func SGR3(n1, n2, n3 int) ColorSequence {
	return ColorSequence(3) |
		(ColorSequence(n1) << colorCodeBitSize) |
		(ColorSequence(n2) << (colorCodeBitSize * 2)) |
		(ColorSequence(n3) << (colorCodeBitSize * 3))
}

func SGR4(n1, n2, n3, n4 int) ColorSequence {
	return ColorSequence(4) |
		(ColorSequence(n1) << colorCodeBitSize) |
		(ColorSequence(n2) << (colorCodeBitSize * 2)) |
		(ColorSequence(n3) << (colorCodeBitSize * 3)) |
		(ColorSequence(n4) << (colorCodeBitSize * 4))
}

func (c ColorSequence) WriteTo(w io.Writer) (int64, error) {
	if c <= 0 {
		return 0, nil
	}
	n := int64(0)
	io.WriteString(w, "\x1B[")
	count := c & colorCodeMask
	for {
		c >>= colorCodeBitSize
		_n, err := fmt.Fprintf(w, "%d", c&colorCodeMask)
		n += int64(_n)
		if err != nil {
			break
		}
		count--
		if count <= 0 {
			break
		}
		_n, err = w.Write([]byte{';'})
		n += int64(_n)
		if err != nil {
			break
		}
	}
	_n, err := w.Write([]byte{'m'})
	n += int64(_n)
	return n, err
}

func (B *Buffer) Write(b []byte) (int, error) {
	return B.Out.Write(b)
}

func (B *Buffer) puts(s []Cell) _Range {
	defaultColor := ColorSequence(B.RefreshColor())
	color := ColorSequence(-1)
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
