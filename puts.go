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

type _Range []cellT

func (B *Buffer) puts(s []cellT) _Range {
	B.RefreshColor()
	var color ColorCode
	for _, ch := range s {
		if ch.color != color {
			color = ch.color
			fmt.Fprintf(B.Out, "\x1B[%d;1m", color)
		}
		ch.Moji.PrintTo(B.Out)
	}
	if color != 37 {
		io.WriteString(B.Out, "\x1B[0;37m")
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Moji.Width()
	}
	return
}
