package readline

import (
	"fmt"
	"io"
)

type Coloring interface {
	// Reset has to initialize receiver's fields and return default color.
	Init() ColorSequence
	// Next has to return color for the given rune.
	Next(rune) ColorSequence
}

type ColorSequence int64

const (
	colorCodeBitSize = 8
	colorCodeMask    = (1<<colorCodeBitSize - 1)
)

const (
	Black ColorSequence = 3 | ((30 + iota) << colorCodeBitSize) | (49 << (colorCodeBitSize * 2)) | (1 << (colorCodeBitSize * 3))
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	_
	DefaultForeGroundColor
)

const (
	DarkGray ColorSequence = 3 | ((30 + iota) << colorCodeBitSize) | (22 << (colorCodeBitSize * 2)) | (49 << (colorCodeBitSize * 3))
	DarkRed
	DarkGree
	DarkYellow
	DarkBlue
	DarkMagenta
	DarkCyan
	DarkWhite
)

type _MonoChrome struct{}

func (_MonoChrome) Init() ColorSequence {
	return DefaultForeGroundColor
}

func (_MonoChrome) Next(rune) ColorSequence {
	return DefaultForeGroundColor
}

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
