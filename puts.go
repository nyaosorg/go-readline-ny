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

type _Range []Moji

func (B *Buffer) puts(s []Moji) _Range {
	for _, ch := range s {
		ch.PrintTo(B.Out)
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Width()
	}
	return
}
