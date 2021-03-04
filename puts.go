package readline

import (
	"fmt"
	"io"
)

func (this *Buffer) backspace(n WidthT) {
	if n > 1 {
		fmt.Fprintf(this.Out, "\x1B[%dD", n)
	} else if n == 1 {
		this.Out.WriteByte('\b')
	}
}

func (this *Buffer) eraseline() {
	io.WriteString(this.Out, "\x1B[0K")
}

type _Range []Moji

func (this *Buffer) puts(s []Moji) _Range {
	for _, ch := range s {
		ch.PrintTo(this.Out)
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Width()
	}
	return
}
