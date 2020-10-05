package readline

import (
	"fmt"
	"io"
	"os"
)

var SurrogatePairOk = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != ""

func (this *Buffer) putRune(m Moji) {
	m.Put(this.Out)
}

func (this *Buffer) putRunes(ch Moji, n width_t) {
	if n <= 0 {
		return
	}
	this.putRune(ch)
	for i := width_t(1); i < n; i++ {
		this.putRune(ch)
	}
}

func (this *Buffer) backspace(n width_t) {
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

func (s Range) Width() (w width_t) {
	for _, ch := range s {
		w += ch.Width()
	}
	return
}
