package readline

import (
	"fmt"
	"io"
)

func (this *Buffer) putRune(m Moji) {
	m.PrintTo(this.Out)
}

func (this *Buffer) backspace(n WidthT) {
	if n > 1 {
		fmt.Fprintf(this.Out, "\x1B[%dD", n)
	} else if n == 1 {
		this.Out.WriteByte('\b')
	}
}

func (this *Buffer) Eraseline() {
	io.WriteString(this.Out, "\x1B[0K")
}

type Range []Moji

func (this *Buffer) puts(s []Moji) Range {
	for _, ch := range s {
		this.putRune(ch)
	}
	return Range(s)
}

func (s Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Width()
	}
	return
}
