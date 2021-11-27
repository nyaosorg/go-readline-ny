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

func (B *Buffer) puts(s []_Cell) _Range {
	B.RefreshColor()
	color := _PackedColorCode(White)
	for _, ch := range s {
		if ch.color != color {
			color = ch.color
			fmt.Fprintf(B.Out, "\x1B[%d;1m", color)
		}
		ch.Moji.PrintTo(B.Out)
	}
	if color != White {
		io.WriteString(B.Out, "\x1B[37;1m")
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Moji.Width()
	}
	return
}
